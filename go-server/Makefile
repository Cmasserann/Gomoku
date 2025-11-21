SRC_FILE := server.go \
			struct.go 


BIN_DIR := ./bin


all: build

build:
	go build -o ${BIN_DIR}/server ${SRCS_FILE}

test_run:
	@go run ${SRC_FILE}

run: build
	${BIN_DIR}/server

clean:
	rm -rf ${BIN_DIR}/*

remake:
	fclean all

.PHONY: all build test_run run clean