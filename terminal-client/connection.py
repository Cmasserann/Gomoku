"""All Connection-specific terminal client code."""

import socket

def connect_to_server(address: str, port: int) -> str:
    """Connect to the Gomoku server."""

    HOST = address
    PORT = port
    MESSAGE = b'Hello, World!\n'

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        s.sendall(MESSAGE)
        data = s.recv(1024)

        data_board = data.decode('utf-8').strip()
    return data_board