#!/usr/bin/make -f

VERSION := $(shell git describe)

test: fmt
	go test -count=1 -short -cover $(ARGS) ./...

fmt:
	go mod tidy
	go fmt ./...

build:
	go build ./...

install: test
	go install -ldflags="-X 'main.Version=$(VERSION)'" github.com/mdwhatcott/gtd/v3/cmd/gtd
