import curses

import client_tool

def draw_menu(stdscr: curses.window):
    key = 0
    selection = 0
    stdscr.clear()    
    stdscr.refresh()

    curses.curs_set(0)

    curses.mousemask(curses.ALL_MOUSE_EVENTS)

    curses.start_color()
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)

    title = "Welcome to the Go Game!"
    quit_msg = "Quit. [q]"
    play_msg = "Play the game. [p]"

    while (key != ord('q')):
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if key == curses.KEY_DOWN or key == ord('s'):
            selection = (selection + 1) % 2
        elif key == curses.KEY_UP or key == ord('w'):
            selection = (selection - 1) % 2

        if key == ord('\n'):
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

        if x_title < 0 or x_play < 0 or x_quit < 0 or y_title < 0 or y_play <= 2 or y_quit <= 4:
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

        if key == ord('p'):
            draw_game(stdscr)

def draw_game(stdscr: curses.window):
    key = 0
    turn_to_play = True
    space_pressed = False
    cursor_x = 0
    cursor_y = 0
    x_input = -1
    y_input = -1
    sug_x = -1
    sug_y = -1
    resp = 0

    board = client_tool.get_board()
    if not board:
        stdscr.clear()
        stdscr.addstr(0, 0, "Failed to connect to server. Press any key to exit.")
        stdscr.getch()
        stdscr.clear()
        stdscr.refresh()
        return

    while (key != ord('q')):
        stdscr.clear()
        height, width = stdscr.getmaxyx()

        if height < len(board["board"]) + 7 or width < len(board["board"]) * 2:
            stdscr.addstr(0, 0, "Terminal too small!", curses.color_pair(2))
            stdscr.refresh()
            key = stdscr.getch()
            continue

        if not turn_to_play:
            new_board = client_tool.wait_for_change(board)
            if new_board["board"] != board["board"]:
                board = new_board
                turn_to_play = True

        goban = board["board"]

        for y in range(len(goban)):
            for x in range(len(goban[y])):
                char = "."
                if (x, y) == (cursor_x, cursor_y):
                    stdscr.attron(curses.color_pair(3))
                if goban[y][x] == 1:
                    char = "B"
                elif goban[y][x] == 2:
                    char = "W"
                stdscr.addstr(y, x * 2, char)
                stdscr.attroff(curses.color_pair(3))
        
        stdscr.attron(curses.color_pair(3))
        stdscr.addstr(height - 1, 0, "Press 'q' to exit.")
        stdscr.attroff(curses.color_pair(3))

        if key == curses.KEY_UP:
            cursor_y = max(0, cursor_y - 1)
        elif key == curses.KEY_DOWN:
            cursor_y = min(len(goban) - 1, cursor_y + 1)
        elif key == curses.KEY_LEFT:
            cursor_x = max(0, cursor_x - 1)
        elif key == curses.KEY_RIGHT:
            cursor_x = min(len(goban) - 1, cursor_x + 1)

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

        stdscr.addstr(len(goban) + 2, 0, f"Cursor at x={cursor_x}, y={cursor_y}", curses.color_pair(1))

        if turn_to_play:
            stdscr.addstr(len(goban) + 4, 0, "Your turn to play.", curses.color_pair(1))
        else:
            stdscr.addstr(len(goban) + 4, 0, "Waiting for opponent...", curses.color_pair(1))

        if key in range(ord('0'), ord('9') + 1):
            digit = key - ord('0')
            if not space_pressed :
                if x_input == -1 :
                    x_input = digit
                else :
                    x_input = x_input * 10 + digit
                    if x_input >= len(goban) :
                        x_input = -1
            else :
                if y_input == -1 :
                    y_input = digit
                else :
                    y_input = y_input * 10 + digit
                    if y_input >= len(goban) :
                        y_input = -1
        
        if not space_pressed and key == ord(' ') and x_input != -1 :
            space_pressed = True
        elif key == ord('\n') and turn_to_play :
            if space_pressed and y_input != -1 :
                resp = client_tool.send_move(x_input, y_input, 1)
            elif cursor_x != -1 and cursor_y != -1 :
                resp = client_tool.send_move(cursor_x, cursor_y, 1)
            if resp :
                turn_to_play = False
            x_input = -1
            y_input = -1
            space_pressed = False

        if x_input != -1 :
            stdscr.addstr(len(goban) + 5, 0, f"X input: {x_input}", curses.color_pair(1))
        if y_input != -1 :
            stdscr.addstr(len(goban) + 6, 0, f"Y input: {y_input}", curses.color_pair(1))

        if key == ord('c') :
            x_input = -1
            y_input = -1
            space_pressed = False

        if key == ord('h') :
            ret = client_tool.ai_suggest()
            if ret :
                sug_x = ret.get("x", -1)
                sug_y = ret.get("y", -1)
        if sug_x != -1 and sug_y != -1 :
            stdscr.addstr(len(goban) + 7, 0, f"AI suggests move at x={sug_x}, y={sug_y}", curses.color_pair(1))


        stdscr.refresh()

        stdscr.timeout(200)
        key = stdscr.getch()

    stdscr.timeout(-1)
    stdscr.clear()
    stdscr.refresh()


def main():
    curses.wrapper(draw_menu)

if __name__ == "__main__":
    main()