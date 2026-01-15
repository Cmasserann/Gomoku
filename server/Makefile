SRCS := main.go API.go struct.go AI.go AI_tools.go

BIN_DIR := ./bin

TMP_FILE := ./cover.out \
			./cover.html

NAME := server

all: build run

run: 
	${BIN_DIR}/server

$(NAME): build

build:
	go build -o ${BIN_DIR}/server ${SRCS}

build_inspect:
	go build -x -work -o ${BIN_DIR}/server ${SRCS}

clean:
	rm -rf ${TMP_FILE}

fclean: clean
	rm -rf ${BIN_DIR}


remake:
	fclean all

test:
	go test ./...

test_coverage:
	go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

.PHONY: all build build_inspect test run clean