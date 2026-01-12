import curses

import client_tool

turn_to_play = True
space_pressed = False
x_input = -1
y_input = -1
sug_x = -1
sug_y = -1


def draw_menu(stdscr: curses.window):
    key = 0
    selection = 0
    stdscr.clear()
    stdscr.refresh()

    curses.curs_set(0)

    curses.mousemask(curses.ALL_MOUSE_EVENTS)

    curses.start_color()
    # base color
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    # suggestion color
    curses.init_pair(2, curses.COLOR_BLACK, curses.COLOR_RED)
    # cursor selection color
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)
    # cursor on top of suggestion
    curses.init_pair(4, curses.COLOR_RED, curses.COLOR_WHITE)
    # opponents color
    curses.init_pair(5, curses.COLOR_RED, curses.COLOR_BLACK)
    # oppoenents color cursor
    curses.init_pair(6, curses.COLOR_RED, curses.COLOR_WHITE)
    # your color cursor
    curses.init_pair(7, curses.COLOR_CYAN, curses.COLOR_WHITE)

    title = "Welcome to the Go Game!"
    quit_msg = "Quit. [q]"
    play_msg = "Play the game. [p]"

    while key != ord("q"):
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if key == curses.KEY_DOWN or key == ord("s"):
            selection = (selection + 1) % 2
        elif key == curses.KEY_UP or key == ord("w"):
            selection = (selection - 1) % 2

        if key == ord("\n"):
            if selection == 0:
                draw_game(stdscr)
            else:
                break

        x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
        y_title = int((height // 2) - 2)
        x_play = int((width // 2) - (len(play_msg) // 2) - len(play_msg) % 2)
        y_play = y_title + 2
        x_quit = int((width // 2) - (len(quit_msg) // 2) - len(quit_msg) % 2)
        y_quit = y_play + 2

        if (
            x_title < 0
            or x_play < 0
            or x_quit < 0
            or y_title < 0
            or y_play <= 2
            or y_quit <= 4
        ):
            stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
            stdscr.refresh()
            key = stdscr.getch()
            continue

        stdscr.addstr(y_title, x_title, title)

        if selection == 0:
            stdscr.attron(curses.color_pair(3))
            stdscr.addstr(y_play, x_play, play_msg)
            stdscr.attroff(curses.color_pair(3))
            stdscr.addstr(y_quit, x_quit, quit_msg)
        else:
            stdscr.addstr(y_play, x_play, play_msg)
            stdscr.attron(curses.color_pair(3))
            stdscr.addstr(y_quit, x_quit, quit_msg)
            stdscr.attroff(curses.color_pair(3))

        stdscr.refresh()

        key = stdscr.getch()

        if key == ord("p"):
            draw_game(stdscr)


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

        if height < len(board["board"]) + 5 or width < len(board["board"]) * 2:
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
            draw_big_goban(stdscr, goban, cursor_x, cursor_y, width)
        else:
            draw_goban(stdscr, goban, cursor_x, cursor_y, width)

        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(height - 1, 0, "Press 'q' to exit.")
        stdscr.attroff(curses.color_pair(3))

        cursor_x, cursor_y = client_tool.get_cursor_pos(
            key, cursor_x, cursor_y, len(goban)
        )

        if key == curses.KEY_MOUSE:
            _, mx, my, _, _ = curses.getmouse()
            if big_goban:
                grid_x = mx // 4
                grid_y = (my - 1) // 2
            else:
                grid_x = mx // 2
                grid_y = my
            if grid_x < len(goban) and grid_y < len(goban):
                cursor_x = grid_x
                cursor_y = grid_y
                send_move(grid_x, grid_y, 1)

        start_line = len(goban) * 2 if big_goban else len(goban)
        draw_info_panel(stdscr, start_line  , turn_to_play, big_goban)

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


def draw_goban(
    stdscr: curses.window,
    goban: list[list[int]],
    cursor_x: int,
    cursor_y: int,
    width: int,
):
    start_y = int((width // 2) - (len(goban) // 2) - len(goban) % 2)
    start_y = 0

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
            stdscr.addstr(int(start_y + y), x * 2, char)
            stdscr.attroff(curses.color_pair(3))


def draw_big_goban(
    stdscr: curses.window,
    goban: list[list[int]],
    cursor_x: int,
    cursor_y: int,
    width: int,
):
    start_y = int((width // 2) - (len(goban * 2) // 2) - len(goban) % 2)
    start_y = 0

    for y in range(len(goban)):
        stdscr.addstr(y * 2, 0, "║")
        stdscr.addstr(y * 2 + 1, 0, "║")
        stdscr.addstr(y * 2, len(goban[y] * 4) + 1, "║")
        stdscr.addstr(y * 2 + 1, len(goban[y] * 4) + 1, "║")
        if y == 0:
            stdscr.addstr(y, 0, "╔" + "════" * (len(goban[y]) - 1) + "════╗")
        if y == len(goban) - 1:
            stdscr.addstr(y * 2 + 2, 0, "║")
            stdscr.addstr(y * 2 + 2, len(goban[y] * 4) + 1, "║")
            stdscr.addstr(y * 2 + 3, 0, "╠" + "════" * (len(goban[y]) - 1) + "════╝")
        for x in range(len(goban[y])):
            char: list[str] = ["  ", "__"]
            if goban[y][x] == 1:
                char = ["╔╗", "╚╝"]
            elif goban[y][x] == 2:
                char = ["┌┐", "└┘"]

            color = select_color(x, y, cursor_x, cursor_y, goban[y][x])
            stdscr.attron(curses.color_pair(color))

            # stdscr.addstr(y * 2 + 1, x * 4 + 2, char[0])
            # stdscr.addstr(y * 2 + 2, x * 4 + 2, char[1])
            stdscr.addstr(start_y + y * 2 + 1, x * 4 + 2, char[0])
            stdscr.addstr(start_y + y * 2 + 2, x * 4 + 2, char[1])

            stdscr.attroff(curses.color_pair(3))

def draw_info_panel(stdscr: curses.window, start_line: int, turn_to_play: bool, big_goban: bool):
    if turn_to_play:
        msg = " Your turn to play. "
    else:
        msg = " Waiting for opponent... "
    capt_msg_1 = " B Capture: " + "1 "
    capt_msg_2 = " W Capture: " + "0 "

    if big_goban:
        start_line += 2
        stdscr.addstr(start_line, 0, "║")
        stdscr.addstr(start_line + 1, 0, "╚")

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line, 0 + 1, msg)
        stdscr.attroff(curses.color_pair(1))
        stdscr.addstr(start_line + 1, 0 + 1, "═" * (len(msg) + len(capt_msg_1) + len(capt_msg_2) + 2))

        stdscr.addstr(start_line, len(msg) + 1, "║")
        stdscr.addstr(start_line - 1, len(msg) + 1, "╦")
        stdscr.addstr(start_line + 1, len(msg) + 1, "╩")

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line, len(msg) + 2, capt_msg_1)
        stdscr.attron(curses.color_pair(5))
        stdscr.addstr(start_line, len(msg) + len(capt_msg_1) + 3, capt_msg_2)
        stdscr.attroff(curses.color_pair(5))

        stdscr.addstr(start_line - 1, len(msg) + len(capt_msg_1) + 2, "╦")
        stdscr.addstr(start_line, len(msg) + len(capt_msg_1) + 2, "║")
        stdscr.addstr(start_line + 1, len(msg) + len(capt_msg_1) + 2, "╩")
        stdscr.addstr(start_line - 1, len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3, "╦")
        stdscr.addstr(start_line, len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3, "║")
        stdscr.addstr(start_line + 1, len(msg) + len(capt_msg_1) + len(capt_msg_2) + 3, "╝")
    else:
        
        stdscr.addstr(start_line , 0, msg)

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line + 2, 0, capt_msg_1)
        stdscr.attroff(curses.color_pair(1))

        stdscr.attron(curses.color_pair(5))
        stdscr.addstr(start_line + 3, 0, capt_msg_2)
        stdscr.attroff(curses.color_pair(5))



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


def main():
    curses.wrapper(draw_menu)


if __name__ == "__main__":
    main()
