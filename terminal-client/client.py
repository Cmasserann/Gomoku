"""Terminal client for Gomoku game."""
# Connect to the server, send x and y, receive updates, and display using curses.

# import sys
# import os
import curses

import tcurses
import connection


def main():
    stdscr = tcurses.setup_curses()
    windows = "h"
    k = 0

    try:
        while (k != ord('q')):
            if k == ord('i'):
                windows = "i"
            elif k == ord('h'):
                windows = "h"
                

            if windows == "h":
                tcurses.welcome_Screen(stdscr, k)

            if windows == "i":
                data = connection.connect_to_server("127.0.0.1", 8080)

                print(f"Data from server: {data}")
                toboard = tcurses.str_to_data(data)
                print(f"to board: {toboard}")
                board = tcurses.get_board(toboard)
                tcurses.update_screen(stdscr, board)
                windows = ""

                # tcurses.update_screen(stdscr, board)
                # stdscr.getch()
            k = stdscr.getch()

    finally:
        curses.nocbreak()
        stdscr.keypad(False)
        curses.echo()
        curses.endwin()


if __name__ == "__main__":
    main()