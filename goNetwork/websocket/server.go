package websocket

import (
	"log"
	"net"
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

type WebSocketConn struct {
	writebuf []byte
	maskKey  [4]byte
	conn     net.Conn
}

func (wc *WebSocketConn) ReadData() (data []byte, err error) {

	var tmpData []byte // store shard datas

	var mask []byte

	var index int

	var b []byte // each data

	for {
		b = make([]byte, 1024)
		// read first two bytes of the frame
		if _, err := wc.conn.Read(b[:2]); err != nil {
			return nil, err
		}

		final := b[0]&finalBit != 0 // whether a final frame
		dataType := int(b[0] & 0xf) // frame type
		index = index + 1

		payloadLen := b[1] & 0x7f // count data length
		// dataLen := int64(payloadLen)
		if payloadLen == 126 {
			// dataLen = int64(binary.BigEndian.Uint16(b[:2]))
			index = index + 2
		} else if payloadLen == 127 {
			// dataLen = int64(binary.BigEndian.Uint64(b[:8]))
			if _, err := wc.conn.Read(b[:8]); err != nil {
				return nil, err
			}
			index = index + 6
		}
		log.Printf("payload length %d\r\n", payloadLen)

		if b[1]&maskBit != 0 { // has mask encode
			if _, err := wc.conn.Read(wc.maskKey[:]); err != nil {
				return nil, err
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
				return
			case 9:
				data = append([]byte{0x8A, byte(len(data))}, data...)
				return
			case 0, 1, 2:
				data = append(tmpData, data...)
				return
			}
		} else {
			tmpData = append(tmpData, data...)
		}
	}
}

func maskByte(key [4]byte, data []byte, index int) {
	for i := range data {
		data[i] ^= key[index&3]
		index++
	}
}
