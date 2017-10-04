all: build test

init:
	go get -u github.com/prometheus/promu
	go get -u github.com/AlekSi/gocoverutil

build:
	go install -v
	promu build

test:
	go test -v -race
	gocoverutil -coverprofile=coverage.txt test -v
