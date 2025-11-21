hello:
	echo "Hello"

build:
	go build -o bin/main test.go

run:
	go run main.go