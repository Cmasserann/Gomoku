import curses

import game

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
            break

    if key == ord("p") or (key == ord("\n") and selection == 0):
        game.draw_game(stdscr)