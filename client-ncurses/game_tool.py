import curses
from typing import Any
import requests  # type: ignore
import time

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


def wait_for_change(old_board: dict[str, Any]):
    while True:
        new_board = get_board()
        if new_board["board"] != old_board["board"]:
            return new_board
        time.sleep(0.2)


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


def draw_info_panel(
    stdscr: curses.window,
    start_line: int,
    start_x: int,
    turn_to_play: bool,
    big_goban: bool,
    captures_B: int = 0,
    captures_W: int = 0,
    invite_token: str = "",
):
    if turn_to_play:
        msg = " Your turn to play. "
    else:
        msg = " Waiting for opponent... "
    capt_msg_1 = " B Capture: " + str(captures_B) + " "
    capt_msg_2 = " W Capture: " + str(captures_W) + " "

    if big_goban:
        start_line += 2
        stdscr.addstr(start_line, start_x, "║")
        stdscr.addstr(start_line + 1, start_x, "╚")

        stdscr.addstr(start_line, start_x + 1, msg)
        stdscr.addstr(
            start_line + 1,
            start_x + 1,
            "═" * (len(msg) + len(capt_msg_1) + len(capt_msg_2) + 2),
        )

        stdscr.addstr(start_line, start_x + len(msg) + 1, "║")
        stdscr.addstr(start_line - 1, start_x + len(msg) + 1, "╦")
        stdscr.addstr(start_line + 1, start_x + len(msg) + 1, "╩")

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line, start_x + len(msg) + 2, capt_msg_1)
        stdscr.attron(curses.color_pair(5))
        stdscr.addstr(start_line, start_x + len(msg) + len(capt_msg_1) + 3, capt_msg_2)
        stdscr.attroff(curses.color_pair(5))

        stdscr.addstr(start_line - 1, start_x + len(msg) + len(capt_msg_1) + 2, "╦")
        stdscr.addstr(start_line, start_x + len(msg) + len(capt_msg_1) + 2, "║")
        stdscr.addstr(start_line + 1, start_x + len(msg) + len(capt_msg_1) + 2, "╩")
        stdscr.addstr(
            start_line - 1,
            start_x + len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3,
            "╦",
        )
        stdscr.addstr(
            start_line, start_x + len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3, "║"
        )
        stdscr.addstr(
            start_line + 1,
            start_x + len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3,
            "╝",
        )
    else:
        stdscr.addstr(start_line, start_x, msg)

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line + 2, start_x, capt_msg_1)
        stdscr.attroff(curses.color_pair(1))

        stdscr.attron(curses.color_pair(5))
        stdscr.addstr(start_line + 3, start_x, capt_msg_2)
        stdscr.attroff(curses.color_pair(5))

        if invite_token:
            stdscr.addstr(start_line + 5, start_x, " Invitation Token: " + invite_token + " ")

