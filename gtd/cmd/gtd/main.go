package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
)

func main() {
	// STARTUP:
	//   for each file in projection folder
	//     unmarshal projection into review queue
	//   for each projection in review queue
	//     compare with canonical projection (loaded from event store)
	//     if different, publish (apply and store) events representing the diff

	HANDLER := wireup.BuildHandler(wireup.Requirements{
		IDFunc: func() string { return uuid.New().String() },
		Reader: eventstore.NewReader(readerFunc, decoderFunc),
		Writer: eventstore.NewWriter(encoderFunc, writerFunc),
	})

	COMMAND := &commands.TrackOutcome{Title: "App Finished"}
	HANDLER.Handle(COMMAND)

	fmt.Println("ID: ", COMMAND.Result.ID)
	fmt.Println("Err:", COMMAND.Result.Error)
}

func decoderFunc(_reader io.Reader) storage.Decoder {
	return json.NewDecoder(_reader, events.Registry())
}
func encoderFunc(_writer io.Writer) storage.Encoder {
	return json.NewEncoder(_writer)
}
func readerFunc(_identifier storage.Identifier) io.ReadCloser {
	return fileFunc(_identifier.ID() + ".md")
}
func writerFunc(_identifier storage.Identifier) io.WriteCloser {
	return fileFunc(_identifier.ID() + ".md")
}
func fileFunc(_name string) *os.File {
	PATH := filepath.Join(storagePath, _name)
	FILE, ERR := os.OpenFile(PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ERR != nil {
		panic(ERR)
	}
	return FILE
}

var storagePath = filepath.Join(os.Getenv("HOME"), "Desktop")
