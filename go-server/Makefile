SRC_FILE := server.go \
			struct.go \
			IA.go \


BIN_DIR := ./bin


all: build

build:
	go build -o ${BIN_DIR}/server ${SRCS_FILE}

run: build
	${BIN_DIR}/server

clean:
	rm -rf ${BIN_DIR}

remake:
	fclean all

test:
	go test ./...

test_coverage:
	go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

.PHONY: all build test_run run clean