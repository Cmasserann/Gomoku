SRCS := server.go struct.go IA.go


BIN_DIR := ./bin
all: run

run: build
	${BIN_DIR}/server

build:
	go build -o ${BIN_DIR}/server ${SRCS}

build_inspect:
	go build -x -work -o ${BIN_DIR}/server ${SRCS}

clean:
	rm -rf ${BIN_DIR}	

remake:
	fclean all

test:
	go test ./...

test_coverage:
	go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html

.PHONY: all build build_inspect test run clean