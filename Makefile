#!/usr/bin/make -f

test:
	go test -short $(ARGS) ./...

install:
	go install github.com/mdwhatcott/gtd/cmd/gtd
