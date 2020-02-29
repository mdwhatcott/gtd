#!/usr/bin/make -f

test: build
	go test -count=1 -short $(ARGS) ./...

build:
	go build ./...

install:
	go install github.com/mdwhatcott/gtd/legacy/cmd/gtd
