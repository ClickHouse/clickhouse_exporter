all: build test

init:
	go get -u github.com/AlekSi/gocoverutil

build:
	go install -v
	go build

test:
	go test -v -race
	gocoverutil -coverprofile=coverage.txt test -v
