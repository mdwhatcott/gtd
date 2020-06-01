package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
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

func BuildApplication() *Application {
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

	return &Application{
		handler: coreWireup.BuildOutcomesHandler(REQUIREMENTS),
		editor:  tempfile.NewEditor(),
		reader:  REQUIREMENTS.Reader,
	}
}

type Application struct {
	handler core.Handler
	editor  ui.Editor
	reader  joyride.StorageReader
}

func (this *Application) editOutcomes(_ids []string) {
	waiter := new(sync.WaitGroup)
	waiter.Add(len(_ids))
	for _, ID := range _ids {
		go this.editOutcome(ID, waiter)
	}
	waiter.Wait()
}

func (this *Application) editOutcome(ID string, waiter *sync.WaitGroup) {
	defer waiter.Done()

	log.Println("Replaying outcome details...")
	START := time.Now()
	STREAM := &storage.OutcomeEventStream{OutcomeID: ID}
	this.reader.Read(STREAM)
	if len(STREAM.Result.Events) == 0 {
		log.Println("No history found for requested outcome ID:", ID)
		return
	}
	PROJECTOR := projections.NewOutcomeDetailsProjector()
	PROJECTOR.Apply(STREAM.Result.Events...)
	log.Printf("Outcome details replayed for id %s in: %s", ID, time.Since(START))
	PROJECTION := PROJECTOR.OutcomeDetailsProjection()
	RESULT := this.editor.EditTempFile(ux.FormatOutcomeDetails(PROJECTION))
	PARSER := ux.NewOutcomeDetailParser(this.handler, ID, PROJECTION, RESULT)
	ERR := PARSER.Parse()
	if ERR != nil {
		log.Fatal(ERR)
	}
}

func (this *Application) PresentOutcomesListing(contexts []string) {
	PROJECTION := replayOutcomesListing(this.reader)
	RESULT := this.editor.EditTempFile(ux.FormatOutcomesListing(PROJECTION))
	PARSER := ux.NewOutcomesListingParser(this.handler, PROJECTION, RESULT)
	this.editOutcomes(PARSER.Parse())
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

func (this *Application) PresentIncompleteActions(contexts []string) {
	PROJECTION := replayIncompleteActions(this.reader)
	RESULT := this.editor.EditTempFile(ux.FormatIncompleteActions(PROJECTION))
	PARSER := ux.NewIncompleteActionsParser(this.handler, PROJECTION, RESULT)
	this.editOutcomes(PARSER.Parse())
}

func replayIncompleteActions(_reader joyride.StorageReader) projections.IncompleteActionsByContext {
	log.Println("Replaying incomplete actions...")
	START := time.Now()
	STREAM := &storage.EventStream{}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewIncompleteActionsByContextProjector()
	PROJECTOR.Apply(STREAM.Result.Events...)
	log.Println("Incomplete actions replayed in:", time.Since(START))
	return PROJECTOR.IncompleteActionsByContextProjection()
}

func (this *Application) PresentIncompleteAction(contexts []string) {
	// TODO
}

func (this *Application) PresentContexts() {
	PROJECTION := replayIncompleteActions(this.reader)
	for _, context := range PROJECTION.Contexts {
		fmt.Printf("- %s (%d)\n", context.Name, len(context.Actions))
	}
}
