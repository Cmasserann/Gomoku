# 19*19
# {[{3 3} {16 17}] [{16 16}]}


def print_board(board: tuple[list[tuple[int, int]], list[tuple[int, int]]]) -> None:
    size = 19

    for y in range(size):
        for x in range(size):
            if (x, y) in board[0]:
                print("B", end=" ")
            elif (x, y) in board[1]:
                print("W", end=" ")
            else:
                print(".", end=" ")
        print()


if __name__ == "__main__":
    print_board(([(3, 3), (16, 17)], [(16, 16)]))