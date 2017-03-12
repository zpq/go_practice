package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"net/http"
)

const (
	finalBit      = 1 << 7
	maskBit       = 1 << 7
	textMessage   = 1
	binaryMessage = 2
	closeMessaage = 8
	pingMessage   = 9
	pongMessage   = 10
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

type WebSocketConn struct {
	writeBuf []byte
	maskKey  [4]byte
	conn     net.Conn
}

func NewConn(conn net.Conn) *WebSocketConn {
	return &WebSocketConn{conn: conn}
}

func (wc *WebSocketConn) Close() {
	wc.conn.Close()
}

func (wc *WebSocketConn) SendData(data []byte, messageType int) {
	length := len(data)
	wc.writeBuf = make([]byte, 10+length)
	// 数据开始和结束的位置
	payloadStart := 2

	// 数据帧的第一个字节, 不支持分片，且值能发送文本类型数据
	// 所以二进制位为 %b1000 0001
	// b0 := []byte{0x81}
	wc.writeBuf[0] = byte(messageType) | finalBit

	// 数据帧第二个字节，服务器发送的数据不需要进行掩码处理
	switch {
	case length >= 65536:
		wc.writeBuf[1] = byte(0x00) | 127
		binary.BigEndian.PutUint64(wc.writeBuf[payloadStart:], uint64(length))
		// 需要 8 byte 来存储数据长度
		payloadStart += 8
	case length > 125:
		wc.writeBuf[1] = byte(0x00) | 126
		binary.BigEndian.PutUint16(wc.writeBuf[payloadStart:], uint16(length))
		// 需要 2 byte 来存储数据长度
		payloadStart += 2
	default:
		wc.writeBuf[1] = byte(0x00) | byte(length)
	}
	copy(wc.writeBuf[payloadStart:], data[:])
	wc.conn.Write(wc.writeBuf[:payloadStart+length])
}

func (wc *WebSocketConn) ReadData() (data []byte, messageType int, err error) {

	var tmpData []byte // store shard datas

	var index int

	var b []byte // each data

	for {
		b = make([]byte, 4096)
		// read first two bytes of the frame
		if _, err := wc.conn.Read(b[:2]); err != nil {
			return nil, 0, err
		}

		final := b[0]&finalBit != 0 // whether a final frame
		log.Println("isfinal: ", final)
		dataType := int(b[0] & 0xf) // frame type
		switch dataType {
		case textMessage:
			messageType = textMessage
		case binaryMessage:
			messageType = binaryMessage
		}
		log.Println("messageType: ", messageType)
		index = index + 1

		payloadLen := b[1] & 0x7f // count data length
		// dataLen := int64(payloadLen)
		if payloadLen == 126 {
			// dataLen = int64(binary.BigEndian.Uint16(b[:2]))
			index = index + 2
		} else if payloadLen == 127 {
			// dataLen = int64(binary.BigEndian.Uint64(b[:8]))
			if _, err := wc.conn.Read(b[:8]); err != nil {
				return nil, 0, err
			}
			index = index + 6
		}
		log.Printf("payload length %d\r\n", payloadLen)

		if b[1]&maskBit != 0 { // has mask encode
			if _, err := wc.conn.Read(wc.maskKey[:]); err != nil {
				return nil, 0, err
			}
			index = index + 4
		}

		if len(wc.maskKey) > 0 {
			maskByte(wc.maskKey, b, index)
		}

		data = b[index:]
		if final {
			switch dataType {
			case 8: //close frame
				wc.conn.Close()
				data = []byte{0x88, 0x00}
				return data, 0, errors.New("client close connection")
			case 9:
				data = append([]byte{0x8A, byte(len(data))}, data...)
				return
			case 0, 1, 2:
				data = append(tmpData, data...)
				return
			}
		} else {
			tmpData = append(tmpData, data...)

			log.Println("data: ", string(tmpData))
		}
	}
}

func Upgrade(w http.ResponseWriter, r *http.Request) (c *WebSocketConn, err error) {
	// 是否是 GET 方法
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return nil, errors.New("websocket: method not GET")
	}
	// 检查 Sec-WebSocket-Version 版本
	if values := r.Header["Sec-Websocket-Version"]; len(values) == 0 || values[0] != "13" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, errors.New("websocket: version != 13")
	}

	// 检查 Connection 和 Upgrade
	if v := r.Header["Connection"]; len(v) == 0 || v[0] != "Upgrade" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, errors.New("websocket: could not find connection header with token 'Upgrade'")
	}
	if v := r.Header["Upgrade"]; len(v) == 0 || v[0] != "websocket" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, errors.New("websocket: could not find connection header with token 'websocket'")
	}

	// 计算 Sec-WebSocket-Accpet 的值
	challengeKey := r.Header.Get("Sec-Websocket-Key")
	if challengeKey == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, errors.New("websocket: key missing or blank")
	}

	var (
		netConn net.Conn
		br      *bufio.Reader
	)

	h, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, errors.New("websocket: response dose not implement http.Hijacker")
	}
	var rw *bufio.ReadWriter
	netConn, rw, err = h.Hijack()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, err
	}
	br = rw.Reader

	if br.Buffered() > 0 {
		netConn.Close()
		return nil, errors.New("websocket: client sent data before handshake is complete")
	}

	// 构造握手成功后返回的 response
	p := []byte{}
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n\r\n"...)

	if _, err = netConn.Write(p); err != nil {
		netConn.Close()
		log.Println("write error! conn close: ", err.Error())
		return nil, err
	}
	log.Println("Upgrade http to websocket successfully")

	return NewConn(netConn), nil
}

func maskByte(key [4]byte, data []byte, index int) {
	for i := range data {
		data[i] ^= key[index&3]
		index++
	}
}

func computeAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
