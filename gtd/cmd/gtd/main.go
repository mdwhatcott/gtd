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
	//   for each file in event store folder
	//     apply events to aggregate to build in-memory repo
	//     apply events to projection to build in-memory repo
	//   for each file in projection folder
	//     unmarshal projection into review queue
	//   for each projection in review queue
	//     compare with canonical projection loaded from event store
	//     if different, publish (apply and store) events representing the diff
	//   All stored aggregates and projections should be up to date w/ actual on disk

	handler := wireup.BuildHandler(wireup.Requirements{
		IDFunc: func() string { return uuid.New().String() },
		Reader: eventstore.NewReader(readerFunc, decoderFunc),
		Writer: eventstore.NewWriter(encoderFunc, writerFunc),
	})

	command := &commands.TrackOutcome{Title: "App Finished"}
	handler.Handle(command)

	fmt.Println("ID: ", command.Result.ID)
	fmt.Println("Err:", command.Result.Error)
}

func decoderFunc(reader io.Reader) storage.Decoder {
	return json.NewDecoder(reader, events.Registry())
}
func encoderFunc(writer io.Writer) storage.Encoder {
	return json.NewEncoder(writer)
}
func readerFunc(identifier storage.Identifier) io.ReadCloser {
	id := identifier.ID()
	path := filepath.Join(storagePath, id+".md")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}
func writerFunc(identifier storage.Identifier) io.WriteCloser {
	id := identifier.ID()
	path := filepath.Join(storagePath, id+".md")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

var storagePath = filepath.Join(os.Getenv("HOME"), "Desktop")