import curses

import menu

def main():
    curses.wrapper(menu.draw_menu)


if __name__ == "__main__":
    main()
