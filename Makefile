nix:
	nix develop --extra-experimental-features "nix-command flakes"

build:
	go build -o bin/main test.go

doc:
	godoc -http=:6060 

run:
	go run main.go

ncurses:
	cd ./client-ncurses && ${MAKE} run

.PHONY: nix build doc run ncurses