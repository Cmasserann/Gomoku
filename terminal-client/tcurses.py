"""All Curses-specific terminal client code."""

import curses

def setup_curses():
    """Set up the curses environment."""
    stdscr = curses.initscr()
    curses.noecho()
    curses.cbreak()
    stdscr.keypad(True)
    curses.start_color()
    curses.init_pair(1, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_BLACK, curses.COLOR_WHITE)
    
    return stdscr

def update_screen(stdscr: curses.window, content: str):
    """Update the screen with the given content."""
    stdscr.clear()
    stdscr.addstr(0, 0, content)
    stdscr.refresh()

def str_to_data(data: str) -> tuple[list[tuple[int, int]], list[tuple[int, int]]]:
    """Convert string representation of the board to internal format."""
    board: tuple[list[tuple[int, int]], list[tuple[int, int]]] = ([], [])

    data = data[2:-2]
    rows = data.split('] [')

    for y, row in enumerate(rows):
        row = row[1:-1]
        cells = row.split('} {')
        for _, cell in enumerate(cells):
            cell = cell.split(' ')
            board[y].append((int(cell[0]), int(cell[1])))

    return board

def get_board(board: tuple[list[tuple[int, int]], list[tuple[int, int]]]) -> str:
    """Convert the board state into a string representation."""
    size = 19
    board_str = ""
    for y in range(size):
        for x in range(size):
            if (x, y) in board[0]:
                board_str += "X "
            elif (x, y) in board[1]:
                board_str += "O "
            else:
                board_str += ". "
        board_str += "\n"
    return board_str

def welcome_Screen(stdscr: curses.window, k: int):
    """Display the welcome screen."""
    stdscr.clear()
    height, width = stdscr.getmaxyx()
    
    title = "Welcome to Gomoku Terminal Client!"[:width-1]
    subtitle = "Written by rvinour and cmassera"[:width-1]
    keystr = "Last key pressed: {}".format(k)[:width-1]
    statusbarstr = "Press 'q' to exit | STATUS BAR"

    # Centering calculations
    start_x_title = int((width // 2) - (len(title) // 2) - len(title) % 2)
    start_x_subtitle = int((width // 2) - (len(subtitle) // 2) - len(subtitle) % 2)
    start_x_keystr = int((width // 2) - (len(keystr) // 2) - len(keystr) % 2)
    start_y = int((height // 2) - 2)

    # Render status bar
    stdscr.attron(curses.color_pair(3))
    stdscr.addstr(height-1, 0, statusbarstr)
    stdscr.addstr(height-1, len(statusbarstr), " " * (width - len(statusbarstr) - 1))
    stdscr.attroff(curses.color_pair(3))

    # Turning on attributes for title
    stdscr.attron(curses.color_pair(2))
    stdscr.attron(curses.A_BOLD)

    # Rendering title
    stdscr.addstr(start_y, start_x_title, title)

    # Turning off attributes for title
    stdscr.attroff(curses.color_pair(2))
    stdscr.attroff(curses.A_BOLD)

    # Print rest of text
    stdscr.addstr(start_y + 1, start_x_subtitle, subtitle)
    stdscr.addstr(start_y + 3, (width // 2) - 2, '-' * 4)
    stdscr.addstr(start_y + 5, start_x_keystr, keystr)

    # Refresh the screen
    stdscr.refresh()