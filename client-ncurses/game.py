import curses

import client_tool
import menu

turn_to_play = True
space_pressed = False
x_input = -1
y_input = -1
sug_x = -1
sug_y = -1


def draw_game(stdscr: curses.window):
    global turn_to_play
    global space_pressed
    global x_input
    global y_input
    global sug_x
    global sug_y

    key = 0
    cursor_x = 0
    cursor_y = 0
    big_goban = True

    board = client_tool.get_board()
    if not board:
        stdscr.clear()
        stdscr.addstr(0, 0, "Failed to connect to server. Press any key to exit.")
        stdscr.getch()
        stdscr.clear()
        stdscr.refresh()
        return

    while key != ord("q"):
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if height < len(board["board"]) + 5 or width < len(board["board"]) * 2 - 1:
            stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
            stdscr.refresh()
            key = stdscr.getch()
            continue

        if (
            width >= len(board["board"]) * 4 + 2
            and height >= len(board["board"]) * 2 + 5
        ):
            big_goban = True
        else:
            big_goban = False

        if not turn_to_play:
            new_board = client_tool.wait_for_change(board)
            if new_board["board"] != board["board"]:
                board = new_board
                turn_to_play = True

        goban: list[list[str]] = board["board"]

        if big_goban:
            start_x = int((width // 2) - (len(goban) * 2) - len(goban) % 2)
            draw_big_goban(stdscr, goban, cursor_x, cursor_y, start_x)
        else:
            start_x = int((width // 2) - len(goban) + 1)
            draw_goban(stdscr, goban, cursor_x, cursor_y, start_x)

        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(height - 1, 0, "Press 'q' to exit.")
        stdscr.attroff(curses.color_pair(3))

        cursor_x, cursor_y = client_tool.get_cursor_pos(
            key, cursor_x, cursor_y, len(goban)
        )

        if key == curses.KEY_MOUSE:
            _, mx, my, _, _ = curses.getmouse()
            if big_goban:
                grid_x = (mx - start_x - 2) // 4
                grid_y = (my - 1) // 2
            else:
                grid_x = (mx - start_x) // 2
                grid_y = my
            if grid_x < len(goban) and grid_y < len(goban):
                cursor_x = grid_x
                cursor_y = grid_y
                send_move(grid_x, grid_y, 1)

        start_line = len(goban) * 2 if big_goban else len(goban)
        client_tool.draw_info_panel(
            stdscr, start_line, start_x, turn_to_play, big_goban
        )

        if key in range(ord("0"), ord("9") + 1):
            digit = key - ord("0")
            if not space_pressed:
                if x_input == -1:
                    x_input = digit
                else:
                    x_input = x_input * 10 + digit
                    if x_input >= len(goban):
                        x_input = -1
            else:
                if y_input == -1:
                    y_input = digit
                else:
                    y_input = y_input * 10 + digit
                    if y_input >= len(goban):
                        y_input = -1

        if not space_pressed and key == ord(" ") and x_input != -1:
            space_pressed = True
        elif space_pressed and key == ord(" "):
            space_pressed = False
            y_input = -1
            x_input = -1
        elif key == ord("\n") and turn_to_play:
            if space_pressed and y_input != -1:
                send_move(x_input, y_input, 1)
            elif cursor_x != -1 and cursor_y != -1:
                send_move(cursor_x, cursor_y, 1)

        if key == ord("c"):
            x_input = -1
            y_input = -1
            space_pressed = False

        if key == ord("h"):
            ret = client_tool.ai_suggest()
            if ret:
                sug_x = ret.get("x", -1)
                sug_y = ret.get("y", -1)
        if sug_x != -1 and sug_y != -1:
            stdscr.addstr(
                start_line + 7,
                0,
                f"AI suggests move at x={sug_x}, y={sug_y}",
                curses.color_pair(1),
            )

        stdscr.refresh()

        stdscr.timeout(200)
        key = stdscr.getch()

    stdscr.timeout(-1)
    stdscr.clear()
    stdscr.refresh()
    menu.draw_menu(stdscr)


def draw_goban(
    stdscr: curses.window,
    goban: list[list[int]],
    cursor_x: int,
    cursor_y: int,
    start_x: int,
):

    for y in range(len(goban)):
        for x in range(len(goban[y])):
            char = "."
            if (x, y) == (cursor_x, cursor_y):
                stdscr.attron(curses.color_pair(3))
            if (x, y) == (sug_x, sug_y):
                if (x, y) == (cursor_x, cursor_y):
                    stdscr.attron(curses.color_pair(4))
                else:
                    stdscr.attron(curses.color_pair(2))
            if goban[y][x] == 1:
                char = "B"
            elif goban[y][x] == 2:
                char = "W"
            stdscr.addstr(y, start_x + x * 2, char)
            stdscr.attroff(curses.color_pair(3))


def draw_big_goban(
    stdscr: curses.window,
    goban: list[list[int]],
    cursor_x: int,
    cursor_y: int,
    start_x: int,
):
    for y in range(len(goban)):
        stdscr.addstr(y * 2, start_x, "║")
        stdscr.addstr(y * 2 + 1, start_x, "║")
        stdscr.addstr(y * 2, start_x + len(goban[y] * 4) + 1, "║")
        stdscr.addstr(y * 2 + 1, start_x + len(goban[y] * 4) + 1, "║")
        if y == 0:
            stdscr.addstr(y, start_x, "╔" + "════" * (len(goban[y]) - 1) + "════╗")
        if y == len(goban) - 1:
            stdscr.addstr(y * 2 + 2, start_x, "║")
            stdscr.addstr(y * 2 + 2, start_x + len(goban[y] * 4) + 1, "║")
            stdscr.addstr(
                y * 2 + 3, start_x, "╠" + "════" * (len(goban[y]) - 1) + "════╝"
            )
        for x in range(len(goban[y])):
            char: list[str] = ["  ", "__"]
            if goban[y][x] == 1:
                char = ["╔╗", "╚╝"]
            elif goban[y][x] == 2:
                char = ["┌┐", "└┘"]

            color = select_color(x, y, cursor_x, cursor_y, goban[y][x])
            stdscr.attron(curses.color_pair(color))

            stdscr.addstr(y * 2 + 1, start_x + x * 4 + 2, char[0])
            stdscr.addstr(y * 2 + 2, start_x + x * 4 + 2, char[1])

            stdscr.attroff(curses.color_pair(3))


def select_color(x: int, y: int, cursor_x: int, cursor_y: int, value: int) -> int:
    global sug_x
    global sug_y

    if (x, y) == (sug_x, sug_y):
        if (x, y) == (cursor_x, cursor_y):
            return 4
        return 2
    elif value == 0:
        if (x, y) == (cursor_x, cursor_y):
            return 3
    elif value == 1:
        if (x, y) == (cursor_x, cursor_y):
            return 7
        return 1
    elif value == 2:
        if (x, y) == (cursor_x, cursor_y):
            return 4
        return 5
    return 0


def send_move(x: int, y: int, color: int):
    resp = client_tool.send_move(x, y, color)

    global x_input
    global y_input
    global sug_x
    global sug_y
    global space_pressed
    global turn_to_play

    x_input = -1
    y_input = -1
    sug_x = -1
    sug_y = -1
    if resp:
        space_pressed = False
        turn_to_play = False
