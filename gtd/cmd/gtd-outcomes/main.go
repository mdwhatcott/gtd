package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/projections"
	coreWireup "github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage"
	storageWireup "github.com/mdwhatcott/gtd/gtd/storage/wireup"
	"github.com/mdwhatcott/gtd/gtd/ui/tempfile"
	"github.com/mdwhatcott/gtd/gtd/ui/ux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable required for resolution of event store file.")
	}

	PATH := filepath.Join(GTDPath, "events.json")
	REQUIREMENTS := coreWireup.Requirements{
		IDFunc: storageWireup.GenerateID,
		Reader: storageWireup.BuildEventStoreReader(PATH),
		Writer: storageWireup.BuildEventStoreWriter(PATH),
	}

	HANDLER := coreWireup.BuildOutcomesHandler(REQUIREMENTS)
	EDITOR := tempfile.NewEditor()

	PROJECTION := replayOutcomesListing(REQUIREMENTS.Reader)
	RESULT := EDITOR.EditTempFile(ux.FormatOutcomesListing(PROJECTION))
	PARSER := ux.NewOutcomesListingParser(HANDLER, PROJECTION, RESULT)
	EDITS := PARSER.Parse()

	for _, EDIT := range EDITS {
		PROJECTION := replayOutcomeDetails(EDIT, REQUIREMENTS.Reader)
		RESULT := EDITOR.EditTempFile(ux.FormatOutcomeDetails(PROJECTION))
		PARSER := ux.NewOutcomeDetailParser(HANDLER, EDIT, PROJECTION, RESULT)
		ERR := PARSER.Parse()
		if ERR != nil {
			log.Fatal(ERR)
		}
	}
}

func replayOutcomesListing(_reader joyride.StorageReader) projections.OutcomesListing {
	STREAM := &storage.EventStream{}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewOutcomesListingProjector()
	PROJECTOR.Apply(STREAM.Result.Events...)
	return PROJECTOR.Projection()
}

func replayOutcomeDetails(_outcomeID string, _reader joyride.StorageReader) projections.OutcomeDetails {
	STREAM := &storage.OutcomeEventStream{OutcomeID: _outcomeID}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewOutcomeDetailsProjector()
	PROJECTOR.Apply(STREAM.Result.Events...)
	return PROJECTOR.Projection()
}
