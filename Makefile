all: reset install test build

build:
	go build -o brainfuck *.go

test:
	go test -v

install:
	go get -d -v
	go install -v
