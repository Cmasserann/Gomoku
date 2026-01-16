NAME := gomoku

nix:
	nix develop --extra-experimental-features "nix-command flakes"

run:
	cd ./server && ${MAKE} run

build:
	cd ./server && ${MAKE} build

all: build run

clean:
	cd ./server && ${MAKE} clean

fclean: clean
	cd ./server && ${MAKE} fclean

ncurses:
	cd ./client-ncurses && ${MAKE} run

re: fclean all

$(NAME) : all

.PHONY: nix run build all clean fclean ncurses