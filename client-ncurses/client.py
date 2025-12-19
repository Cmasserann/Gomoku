# send tcp message to 127.0.0.1:8080

import socket

def send_tcp_message():

    HOST = '127.0.0.1'  # The server's hostname or IP address
    PORT = 8080         # The port used by the server
    MESSAGE = b'Hello, World!\n' # The message to be sent

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        s.sendall(MESSAGE)
        print('Sent', repr(MESSAGE))
        data = s.recv(1024)
        print('Received', repr(data))


if __name__ == "__main__":
    send_tcp_message()