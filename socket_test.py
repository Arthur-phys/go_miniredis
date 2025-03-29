import socket

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('localhost', 8000))

# Socket stays open here! You can send/receive repeatedly.
s.sendall(b"*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n")
print(s.recv(1024))

s.sendall(b"*2\r\n$3\r\nGET\r\n$1\r\nB\r\n")  # Still connected!
print(s.recv(1024))
s.close()  # Explicitly close when done.