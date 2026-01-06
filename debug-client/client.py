import requests
import time

URL_BASE = "http://127.0.0.1:8080"

def get_board():

    try:
        response = requests.get(f"{URL_BASE}/board")
        return response.json()
    except Exception as e:
        print(f"Connection fail : {e}")
        return None

def send_move(x, y, color):

    payload = {"x": x, "y": y, "color": color}
    try:
        response = requests.post(f"{URL_BASE}/move", json=payload)
        if response.status_code != 200:
            print(f"Move rejected: {response.json().get('status', 'Unknown error')}")
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None

def print_board(data):

    if not data or 'board' not in data:
        return
    grid = data['board']
    size = len(grid)
    
    # Header pour les colonnes
    print("   0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18")
    for y in range(size):
        print(f"{y:2} ", end="")
        for x in range(size):
            cell = grid[y][x]
            char = "."
            if cell == 1: char = "B"
            elif cell == 2: char = "W"
            print(f"{char} ", end=" ")
        print()

def wait_for_change(old_board):
    print("AI search the best way to kick your ass...")
    while True:
        new_board = get_board()
        if new_board['board'] != old_board['board']:
            return new_board
        time.sleep(0.2)


def main():

    current_board_data = get_board()
    my_color = 1 

    while True:
        print_board(current_board_data)
        
        try:
            line = input(f"\nJoueur {my_color} - Entrez x y ou 'q' pour quitter : ")
            if line.lower() == 'q':
                break
                
            x_input, y_input = map(int, line.split())
        except ValueError:
            print("Entrée invalide. Merci de taper deux nombres séparés par un espace.")
            continue

        print(f"Envoi du coup ({x_input}, {y_input})...")
        board = send_move(x_input, y_input, my_color)

        if board:
            current_board_data = wait_for_change(board)

if __name__ == "__main__":
    main()
