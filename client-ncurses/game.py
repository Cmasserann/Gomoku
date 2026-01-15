import curses

import game_tool as tool

turn_to_play = True
sug_x = -1
sug_y = -1
token = ""


def draw_game(
    stdscr: curses.window,
    ai_mode: bool = False,
    local_mode: bool = False,
    invite_token: str = "",
):
    global turn_to_play
    global sug_x
    global sug_y
    global token

    key = 0
    cursor_x = 0
    cursor_y = 0
    big_goban = True
    token2 = ""
    turn = 1
    player_2 = False
    win = None

    token, token2 = tool.handle_token(invite_token, ai_mode, local_mode)
    board = tool.get_board()
    if not token or not board:
        msg = "Failed to connect to server."
        if token2 == "1":
            msg = "Invalid invitation token."
        elif token2 == "2":
            msg = "Failed to create a room."
        tool.breaking_error(stdscr, msg)
        return

    if ai_mode is True:
        player_2 = False
        local_mode = False
    elif token2 == "":
        player_2 = True

    while key != ord("q"):
        win = None
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if height < len(board["board"]) + 7 or width < len(board["board"]) * 2 - 1:
            stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
            stdscr.refresh()
            key = stdscr.getch()
            continue

        if (
            width >= len(board["board"]) * 4 + 2
            and height >= len(board["board"]) * 2 + 7
        ):
            big_goban = True
        else:
            big_goban = False

        if not turn_to_play:
            board = tool.get_board()
            if not board:
                stdscr.timeout(-1)
                tool.breaking_error(stdscr, "Connection lost.")
                return

            if board["goban_free"] is True:
                if player_2:
                    win = 2 if turn % 2 == 0 else 1
                else:
                    win = 1 if turn % 2 == 0 else 2
                break
            if turn != board["turn"]:
                board = tool.get_board()
                if board["goban_free"] is True:
                    win = 1 if turn % 2 == 0 else 2
                    break
                turn = board["turn"]
        if local_mode:
            turn_to_play = True
        else:
            if turn % 2 == 1 and not player_2:
                turn_to_play = True
            elif turn % 2 == 0 and player_2:
                turn_to_play = True
            else:
                turn_to_play = False

        goban: list[list[int]] = board["board"]

        if big_goban:
            start_x = int((width // 2) - (len(goban) * 2) - len(goban) % 2)
            draw_big_goban(stdscr, goban, cursor_x, cursor_y, start_x)
        else:
            start_x = int((width // 2) - len(goban) + 1)
            draw_goban(stdscr, goban, cursor_x, cursor_y, start_x)

        stdscr.addstr(height - 1, 0, "Press 'q' to exit.", curses.color_pair(3))

        cursor_x, cursor_y = tool.get_cursor_pos(key, cursor_x, cursor_y, len(goban))

        if key == curses.KEY_MOUSE:
            _, mx, my, _, _ = curses.getmouse()
            if big_goban:
                grid_x = (mx - start_x - 2) // 4
                grid_y = (my - 1) // 2
            else:
                grid_x = (mx - start_x) // 2
                grid_y = my
            if (
                grid_x < len(goban)
                and grid_y < len(goban)
                and grid_x >= 0
                and grid_y >= 0
            ):
                cursor_x = grid_x
                cursor_y = grid_y
                win = send_move(grid_x, grid_y, 1)
                if win:
                    break
                turn += 1

        start_line = len(goban) * 2 if big_goban else len(goban)
        token2 = token2 if not local_mode and not ai_mode else ""
        draw_info_panel(
            stdscr,
            start_line,
            start_x,
            turn_to_play,
            big_goban,
            captures_B=board.get("captured_b", 0),
            captures_W=board.get("captured_w", 0),
            invite_token=token2,
            local_mode=local_mode,
            turn=turn,
        )

        if key == ord("\n") and turn_to_play and cursor_x != -1 and cursor_y != -1:
            move = send_move(cursor_x, cursor_y, 1)
            if move:
                break
            turn += 1

        if key == ord("h"):
            ret = tool.ai_suggest(token)
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
        if key == ord("q"):
            break

        if key == ord("t"):
            tool.debug()
            board = tool.get_board()

        stdscr.refresh()

        stdscr.timeout(200)
        key = stdscr.getch()

    stdscr.timeout(-1)
    stdscr.clear()
    if win is None:
        tool.give_up(token)
    else:
        draw_endGame(stdscr, win)
    stdscr.refresh()


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


def send_move(x: int, y: int, color: int) -> int | None:
    global token

    resp = tool.send_move(x, y, color, token)

    global sug_x
    global sug_y
    global turn_to_play

    sug_x = -1
    sug_y = -1
    if resp:
        turn_to_play = False
        if resp.get("winner") is not None:
            return resp["winner"]
    return None


def draw_info_panel(
    stdscr: curses.window,
    start_line: int,
    start_x: int,
    turn_to_play: bool,
    big_goban: bool,
    captures_B: int = 0,
    captures_W: int = 0,
    invite_token: str = "",
    local_mode: bool = False,
    turn: int = 1,
):
    if turn_to_play:
        msg = " Your turn to play. "
    else:
        msg = " Waiting for opponent... "
    if local_mode:
        msg = (
            " Player 1's turn to play. "
            if turn % 2 == 1
            else " Player 2's turn to play. "
        )
    capt_msg_1 = f" B Capture: {str(captures_B)} "
    capt_msg_2 = f" W Capture: {str(captures_W)} "
    turn_msg = f" Turn: {turn} "
    invite_msg = f" {invite_token} "

    if big_goban:

        def draw_border(
            stdscr: curses.window, start_line: int, start_x: int, last_char: str = "╩"
        ):
            stdscr.addstr(start_line - 1, start_x, "╦")
            stdscr.addstr(start_line, start_x, "║")
            stdscr.addstr(start_line + 1, start_x, last_char)

        start_line += 2
        stdscr.addstr(start_line, start_x, "║")
        stdscr.addstr(start_line + 1, start_x, "╚")

        stdscr.addstr(start_line, start_x + 1, msg)

        line_length = len(msg) + len(capt_msg_1) + len(capt_msg_2) + len(turn_msg) + 3
        if invite_token:
            line_length += len(invite_msg) + 1

        stdscr.addstr(
            start_line + 1,
            start_x + 1,
            "═" * line_length,
        )

        start_x += len(msg) + 1
        draw_border(stdscr, start_line, start_x)
        stdscr.addstr(start_line, start_x + 1, capt_msg_1, curses.color_pair(1))

        start_x += len(capt_msg_1) + 1
        draw_border(stdscr, start_line, start_x)
        stdscr.addstr(start_line, start_x + 1, capt_msg_2, curses.color_pair(5))

        start_x += len(capt_msg_2) + 1
        draw_border(stdscr, start_line, start_x)
        stdscr.addstr(start_line, start_x + 1, turn_msg)

        start_x += len(turn_msg) + 1
        draw_border(stdscr, start_line, start_x, "╝")

        if invite_token:
            stdscr.addstr(start_line + 1, start_x, "╩")
            stdscr.addstr(start_line, start_x + 1, invite_msg)

            start_x += len(invite_msg) + 1
            draw_border(stdscr, start_line, start_x, "╝")

    else:
        stdscr.addstr(start_line, start_x, msg)

        stdscr.addstr(start_line + 1, start_x, turn_msg)

        stdscr.attron(curses.color_pair(1))
        stdscr.addstr(start_line + 2, start_x, capt_msg_1)
        stdscr.attroff(curses.color_pair(1))

        stdscr.attron(curses.color_pair(5))
        stdscr.addstr(start_line + 3, start_x, capt_msg_2)
        stdscr.attroff(curses.color_pair(5))

        if invite_token:
            stdscr.addstr(
                start_line + 5, start_x, " Invitation Token: " + invite_token + " "
            )


def draw_endGame(stdscr: curses.window, winner: int):
    stdscr.clear()
    height, width = stdscr.getmaxyx()

    message = "The Winner is: {}!".format("Player 1" if winner == 1 else "Player 2")

    x_msg = int((width // 2) - (len(message) // 2) - len(message) % 2)
    y_msg = int((height // 2))

    stdscr.addstr(y_msg, x_msg, message)
    stdscr.addstr(y_msg + 2, x_msg, "Press any key to return to menu.")

    stdscr.refresh()
    stdscr.getch()
