from typing import Any
import requests # type: ignore
import time

URL_BASE = "http://127.0.0.1:8080"

def get_board() -> dict[str, Any]:

    try:
        response = requests.get(f"{URL_BASE}/board")
        return response.json()
    except Exception as e:
        print(f"Connection fail : {e}")
        return dict()

def send_move(x: int, y: int, color: int) -> (dict[str, Any] | None):

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
    
def wait_for_change(old_board: dict[str, Any]):
    print("AI search the best way to kick your ass...")
    while True:
        new_board = get_board()
        if new_board['board'] != old_board['board']:
            return new_board
        time.sleep(0.2)

def ai_suggest() -> (dict[str, int] | None):
    try:
        response = requests.get(f"{URL_BASE}/ai-suggest")
        if response.status_code != 200:
            print(f"AI suggestion error: {response.json().get('status', 'Unknown error')}")
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None
