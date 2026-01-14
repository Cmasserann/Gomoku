import curses

import game


def draw_menu(stdscr: curses.window):
    key = 0
    selection = 0
    AI_mode = False
    local_mode = False
    stdscr.clear()
    stdscr.refresh()
    token = ""

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

    while key != ord("q"):
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if key == curses.KEY_DOWN or key == ord("s"):
            selection = (selection + 1) % 5
        elif key == curses.KEY_UP or key == ord("w"):
            selection = (selection - 1) % 5

        if key == ord("\n"):
            if selection == 2:
                local_mode = not local_mode
            elif selection == 3:
                AI_mode = not AI_mode
            else:
                break

        if draw_text(stdscr, width, height, key, selection, AI_mode, local_mode):
            key = stdscr.getch()
            continue

        stdscr.refresh()

        key = stdscr.getch()

        if key == ord("p"):
            break

        if key == ord("a"):
            AI_mode = not AI_mode
        
        if key == ord("l"):
            local_mode = not local_mode
        
        if key == ord("j") or (key == ord("\n") and selection == 1):
            token = draw_join_game(stdscr)
            if token:
                game.draw_game(stdscr, AI_mode, local_mode, token)
                break
            else:
                continue

    if key == ord("p") or (key == ord("\n") and selection == 0):
        game.draw_game(stdscr, AI_mode, local_mode)
    if key == ord("j") or (key == ord("\n") and selection == 1):
        game.draw_game(stdscr, invite_token=token)


def draw_join_game(stdscr: curses.window) -> str:
    curses.echo()
    stdscr.clear()
    while True:
        height, width = stdscr.getmaxyx()

        prompt = "Enter Invitation Token: "
        x_prompt = int((width // 2) - (len(prompt) // 2) - len(prompt) % 2)
        y_prompt = int((height // 2))

        if x_prompt < 0 or y_prompt < 0:
            stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
            stdscr.refresh()

        stdscr.addstr(y_prompt, x_prompt, prompt)
        stdscr.refresh()
        token = stdscr.getstr(y_prompt, x_prompt + len(prompt), 20).decode("utf-8")
        curses.noecho()
        return token

def draw_text(
    stdscr: curses.window,
    width: int,
    height: int,
    key: int,
    selection: int,
    AI_mode: bool,
    local_mode: bool,
) -> bool:
    title = "Welcome to the Go Game!"
    play_msg = "Play the game. [p]"
    join_msg = "Join a room. [j]"
    local_msg = "Local Mode [l]"
    ai_msg = "AI Mode [a]"
    quit_msg = "Quit. [q]"

    x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
    y_title = int((height // 2) - 5)
    x_play = int((width // 2) - (len(play_msg) // 2) - len(play_msg) % 2)
    y_play = y_title + 2
    x_join = int((width // 2) - (len(join_msg) // 2) - len(join_msg) % 2)
    y_join = y_play + 2
    x_local = int((width // 2) - (len(local_msg) // 2) - len(local_msg) % 2)
    y_local = y_join + 2
    x_ai = int((width // 2) - (len(ai_msg) // 2) - len(ai_msg) % 2)
    y_ai = y_local + 2
    x_quit = int((width // 2) - (len(quit_msg) // 2) - len(quit_msg) % 2)
    y_quit = y_ai + 2

    if (
        x_title < 0
        or y_title < 0
        or x_play < 0
        or y_play <= 2
        or x_join < 0
        or y_join <= 4
        or x_local < 0
        or y_local <= 6
        or x_ai < 0
        or y_ai <= 8
        or x_quit < 0
        or y_quit <= 11
    ):
        stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
        stdscr.refresh()
        return True

    stdscr.addstr(y_title, x_title, title)

    if selection == 0:
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(y_play, x_play, play_msg)
        stdscr.attroff(curses.color_pair(3))
    else:
        stdscr.addstr(y_play, x_play, play_msg)

    if selection == 1:
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(y_join, x_join, join_msg)
        stdscr.attroff(curses.color_pair(3))
    else:
        stdscr.addstr(y_join, x_join, join_msg)


    if selection == 2:
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(y_local, x_local, local_msg)
        stdscr.attroff(curses.color_pair(3))
    else:
        stdscr.addstr(y_local, x_local, local_msg)

    if selection == 3:
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(y_ai, x_ai, ai_msg)
        stdscr.attroff(curses.color_pair(3))
    else:
        stdscr.addstr(y_ai, x_ai, ai_msg)

    if selection == 4:
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(y_quit, x_quit, quit_msg)
        stdscr.attroff(curses.color_pair(3))
    else:
        stdscr.addstr(y_quit, x_quit, quit_msg)

    if AI_mode:
        stdscr.addstr(y_ai, x_ai + len(ai_msg) + 1, "[ON]", curses.color_pair(3))
        stdscr.addstr(y_ai, x_ai + len(ai_msg) + 6, "[OFF]")
    else:
        stdscr.addstr(y_ai, x_ai + len(ai_msg) + 1, "[ON]")
        stdscr.addstr(y_ai, x_ai + len(ai_msg) + 6, "[OFF]", curses.color_pair(3))

    if local_mode:
        stdscr.addstr(
            y_local, x_local + len(local_msg) + 1, "[ON]", curses.color_pair(3)
        )
        stdscr.addstr(y_local, x_local + len(local_msg) + 6, "[OFF]")
    else:
        stdscr.addstr(y_local, x_local + len(local_msg) + 1, "[ON]")
        stdscr.addstr(
            y_local, x_local + len(local_msg) + 6, "[OFF]", curses.color_pair(3)
        )
    return False

