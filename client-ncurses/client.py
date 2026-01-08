import curses

import client_tool

def draw_menu(stdscr: curses.window):
    key = 0
    turn_to_play = True

    board = client_tool.get_board()
    if not board:
        stdscr.addstr(0, 0, "Failed to connect to server. Press any key to exit.")
        stdscr.getch()
        return

    stdscr.clear()
    stdscr.refresh()

    curses.mousemask(curses.ALL_MOUSE_EVENTS)

    curses.start_color()
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)

    while (key != ord('q')):
        stdscr.clear()

        if not turn_to_play:
            new_board = client_tool.wait_for_change(board)
            if new_board["board"] != board["board"]:
                board = new_board
                turn_to_play = True

        goban = board["board"]

        for y in range(len(goban)):
            for x in range(len(goban[y])):
                char = "."
                if goban[y][x] == 1:
                    char = "B"
                elif goban[y][x] == 2:
                    char = "W"
                stdscr.addstr(y, x * 2, char)
        
        stdscr.addstr(len(goban) + 1, 0, "Press 'q' to exit.", curses.color_pair(2))

        if key == curses.KEY_MOUSE :
            _, mx, my, _, _ = curses.getmouse()
            stdscr.addstr(len(goban) + 2, 0, f"Mouse clicked at x={mx}, y={my}", curses.color_pair(1))
            grid_x = mx // 2
            grid_y = my
            stdscr.addstr(len(goban) + 3, 0, f"Grid position x={grid_x}, y={grid_y}", curses.color_pair(1))
            if grid_x < len(goban) and grid_y < len(goban):
                resp = client_tool.send_move(grid_x, grid_y, 1)
                if resp:
                    turn_to_play = False


        stdscr.refresh()

        stdscr.timeout(500)
        key = stdscr.getch()



def main():
    curses.wrapper(draw_menu)

if __name__ == "__main__":
    main()