import socket
import time


ADDR = ("127.0.0.1", 10090)
BUF_SIZE = 1024

client = socket.socket()
client.connect(ADDR)

client.send(bytes("hello", encoding="utf-8"))
data = client.recv(BUF_SIZE)
print(data)
client.close()