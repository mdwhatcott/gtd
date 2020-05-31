#!/usr/bin/make -f

test: fmt
	go test -count=1 -short -cover $(ARGS) ./gtd/...

fmt:
	go fmt ./...

build:
	go build ./...

install: test
	go install github.com/mdwhatcott/gtd/gtd/cmd/gtd
