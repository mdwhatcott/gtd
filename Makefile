#!/usr/bin/make -f

test: build
	go test -short $(ARGS) ./...

build:
	go build ./...

install:
	go install github.com/mdwhatcott/gtd/cmd/gtd
