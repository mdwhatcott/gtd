package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/projections"
	coreWireup "github.com/mdwhatcott/gtd/gtd/core/wireup"
	"github.com/mdwhatcott/gtd/gtd/storage"
	storageWireup "github.com/mdwhatcott/gtd/gtd/storage/wireup"
	"github.com/mdwhatcott/gtd/gtd/ui"
	"github.com/mdwhatcott/gtd/gtd/ui/tempfile"
	"github.com/mdwhatcott/gtd/gtd/ui/ux"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

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

	APP := NewApplication(REQUIREMENTS)
	APP.PresentOutcomesListing()
}

type Application struct {
	handler core.Handler
	editor  ui.Editor
	reader  joyride.StorageReader
}

func NewApplication(REQUIREMENTS coreWireup.Requirements) *Application {
	return &Application{
		handler: coreWireup.BuildOutcomesHandler(REQUIREMENTS),
		editor:  tempfile.NewEditor(),
		reader:  REQUIREMENTS.Reader,
	}
}

func (this *Application) PresentOutcomesListing() {
	PROJECTION := replayOutcomesListing(this.reader)
	RESULT := this.editor.EditTempFile(ux.FormatOutcomesListing(PROJECTION))
	PARSER := ux.NewOutcomesListingParser(this.handler, PROJECTION, RESULT)
	this.EditOutcomes(PARSER.Parse())
}

func (this *Application) EditOutcomes(_ids []string) {
	for _, ID := range _ids {
		log.Println("Replaying outcome details...")
		START := time.Now()
		STREAM := &storage.OutcomeEventStream{OutcomeID: ID}
		this.reader.Read(STREAM)
		if len(STREAM.Result.Events) == 0 {
			log.Println("No history found for requested outcome ID:", ID)
			continue
		}
		PROJECTOR := projections.NewOutcomeDetailsProjector()
		PROJECTOR.Apply(STREAM.Result.Events...)
		log.Println("Outcome details replayed in:", time.Since(START))
		PROJECTION := PROJECTOR.OutcomeDetailsProjection()
		RESULT := this.editor.EditTempFile(ux.FormatOutcomeDetails(PROJECTION))
		PARSER := ux.NewOutcomeDetailParser(this.handler, ID, PROJECTION, RESULT)
		ERR := PARSER.Parse()
		if ERR != nil {
			log.Fatal(ERR)
		}
	}

}

func replayOutcomesListing(_reader joyride.StorageReader) projections.OutcomesListing {
	log.Println("Replaying outcomes listing...")
	START := time.Now()
	STREAM := &storage.EventStream{}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewOutcomesListingProjector()
	PROJECTOR.Apply(STREAM.Result.Events...)
	log.Println("Outcomes listing replayed in:", time.Since(START))
	return PROJECTOR.OutcomesListingProjection()
}
