package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
)

func main() {
	requirements := wireup.Requirements{
		IDFunc: func() string { return uuid.New().String() },
		Reader: eventstore.NewReader(reading, decoding),
		Writer: eventstore.NewWriter(encoding, writing),
	}

	// TODO: wrap handler in (cli-)routed controller
	HANDLER := wireup.BuildOutcomesHandler(requirements)
	_ = HANDLER
}

func decoding(_reader io.Reader) storage.Decoder { return json.NewDecoder(_reader, events.Registry()) }
func encoding(_writer io.Writer) storage.Encoder { return json.NewEncoder(_writer) }

func reading() io.ReadCloser  { return open("events.json") }
func writing() io.WriteCloser { return open("events.json") }

func open(_name string) *os.File {
	PATH := filepath.Join(storagePath, _name)
	FILE, ERR := os.OpenFile(PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ERR != nil {
		panic(ERR)
	}
	return FILE
}

var storagePath = filepath.Join(os.Getenv("HOME"), "Desktop")
