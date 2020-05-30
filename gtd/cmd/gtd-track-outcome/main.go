package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
	"github.com/mdwhatcott/gtd/gtd/ui"
	"github.com/mdwhatcott/gtd/gtd/ui/tempfile"
	"github.com/mdwhatcott/gtd/gtd/ui/ux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	HANDLER := wireup.BuildOutcomesHandler(wireup.Requirements{
		IDFunc: func() string { return uuid.New().String() },
		Reader: eventstore.NewReader(reading, decoding),
		Writer: eventstore.NewWriter(encoding, writing),
	})
	CONTENT := tempfile.NewEditor().EditTempFile(ui.TrackOutcomeTemplate)
	PARSER := ux.NewOutcomeDetailParser(HANDLER, "", nil, CONTENT)
	ERR := PARSER.Parse()
	if ERR != nil {
		log.Panic(ERR)
	}
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
