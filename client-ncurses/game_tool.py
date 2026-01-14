import curses
from typing import Any
import requests  # type: ignore

URL_BASE = "http://127.0.0.1:8080"

def create_room(ai_mode: bool, local_mode: bool) -> dict[str, Any] | None:
    payload = {"ai_mode": ai_mode, "local_mode": local_mode}
    try:
        response = requests.post(f"{URL_BASE}/create", json=payload)
        if response.status_code != 200:
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None

def join_room(token:str) -> dict[str, Any] | None:
    payload = {"token": token}
    try:
        response = requests.post(f"{URL_BASE}/join", json=payload)
        if response.status_code != 200:
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None

def get_board() -> dict[str, Any]:
    try:
        response = requests.get(f"{URL_BASE}/board")
        return response.json()
    except Exception as e:
        print(f"Connection fail : {e}")
        return dict()


def send_move(x: int, y: int, color: int, token: str) -> dict[str, Any] | None:
    payload: dict[str, Any] = {"x": x, "y": y, "color": color, "token": token}
    try:
        response = requests.post(f"{URL_BASE}/move", json=payload)
        if response.status_code != 200:
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None


def ai_suggest(token: str) -> dict[str, int] | None:
    payload: dict[str, Any] = {"token": token}
    try:
        response = requests.get(f"{URL_BASE}/ai-suggest", json=payload)
        if response.status_code != 200:
            return None
        return response.json()
    except Exception as e:
        print(f"Error : {e}")
        return None
    
def check_server_online() -> bool:
    try:
        requests.get(f"{URL_BASE}/status")
        return True
    except Exception as e:
        print(f"Connection fail : {e}")
        return False

def give_up(token: str) -> bool:
    payload: dict[str, Any] = {"token": token}
    try:
        response = requests.post(f"{URL_BASE}/giveUp", json=payload)
        if response.status_code != 200:
            return False
        return True
    except Exception as e:
        print(f"Error : {e}")
        return False
    
def debug() -> None:
    payload: dict[str, Any] = {"reset_board": True}
    try:
        response = requests.post(f"{URL_BASE}/debug", json=payload)
        if response.status_code != 200:
            return
        return
    except Exception as e:
        print(f"Error : {e}")
        return


def get_cursor_pos(
    key: int, cursor_x: int, cursor_y: int, goban_size: int
) -> tuple[int, int]:
    if key == curses.KEY_UP:
        cursor_y = max(0, cursor_y - 1)
    elif key == curses.KEY_DOWN:
        cursor_y = min(goban_size - 1, cursor_y + 1)
    elif key == curses.KEY_LEFT:
        cursor_x = max(0, cursor_x - 1)
    elif key == curses.KEY_RIGHT:
        cursor_x = min(goban_size - 1, cursor_x + 1)
    return cursor_x, cursor_y

def handle_token(invite_token: str, ai_mode: bool, local_mode: bool) -> tuple[str, str]:
    if invite_token:
        room = join_room(invite_token)
        if not room:
            return "", "1"
        else:
            return room["token"], ""

    room = create_room(ai_mode, local_mode)
    if not room:
        return "", "2"
    else:
        token = room["player_one"]
        token2 = room["player_two"]
        return token, token2
    
def breaking_error(stdscr: curses.window, msg: str) -> None:
    stdscr.clear()
    stdscr.addstr(0, 0, f"{msg}")
    stdscr.addstr(1, 0, "Press any key to exit.")
    stdscr.refresh()
    stdscr.getch()
    return