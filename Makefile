all: build test

init:
	go get -u github.com/AlekSi/gocoverutil

build:
	go install -v

test:
	docker-compose up -d
	go test -v -race ./...
	gocoverutil -coverprofile=coverage.txt test -v ./...
	docker-compose down
