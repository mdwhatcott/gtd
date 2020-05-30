package main

import (
	"log"
	"os"
	"path/filepath"

	core "github.com/mdwhatcott/gtd/gtd/core/wireup"
	storage "github.com/mdwhatcott/gtd/gtd/storage/wireup"
	"github.com/mdwhatcott/gtd/gtd/ui"
	"github.com/mdwhatcott/gtd/gtd/ui/tempfile"
	"github.com/mdwhatcott/gtd/gtd/ui/ux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gtdPath, ok := os.LookupEnv("GTDPATH")
	if !ok {
		log.Fatal("The 'GTDPATH' environment variable required for resolution of event store file.")
	}

	PATH := filepath.Join(gtdPath, "events.json")
	HANDLER := core.BuildOutcomesHandler(core.Requirements{
		IDFunc: storage.GenerateID,
		Reader: storage.BuildEventStoreReader(PATH),
		Writer: storage.BuildEventStoreWriter(PATH),
	})
	CONTENT := tempfile.NewEditor().EditTempFile(ui.TrackOutcomeTemplate)
	PARSER := ux.NewOutcomeDetailParser(HANDLER, "", nil, CONTENT)
	ERR := PARSER.Parse()
	if ERR != nil {
		log.Panic(ERR)
	}
}
