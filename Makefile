all: build test

init:
	go install github.com/AlekSi/gocoverutil@latest

build:
	go install -v
	go build

test:
	go test -v -race
	gocoverutil -coverprofile=coverage.txt test -v
