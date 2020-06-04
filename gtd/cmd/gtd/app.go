package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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

	PATH := filepath.Join(GTDPath, "events.csv")
	REQUIREMENTS := coreWireup.Requirements{
		IDFunc: storageWireup.GenerateID,
		Reader: storageWireup.BuildCSVEventStoreReader(PATH),
		Writer: storageWireup.BuildCSVEventStoreWriter(PATH),
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

func (this *Application) PresentOutcomesListing(statuses []string) {
	// TODO: filter on provided statuses
	PROJECTION := replayOutcomesListing(this.reader)
	RESULT := this.editor.EditTempFile(ux.FormatOutcomesListing(PROJECTION))
	PARSER := ux.NewOutcomesListingParser(this.handler, PROJECTION, RESULT)
	EDITS := PARSER.Parse()
	this.editOutcomes(EDITS)
	if len(EDITS) > 0 {
		this.PresentOutcomesListing(statuses)
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

func (this *Application) PresentIncompleteActions(contexts []string) {
	PROJECTION := replayIncompleteActions(this.reader)
	CONTEXTS := filterContexts(PROJECTION, contexts)
	RESULT := this.editor.EditTempFile(ux.FormatIncompleteActions(CONTEXTS...))
	PARSER := ux.NewIncompleteActionsParser(this.handler, RESULT, CONTEXTS...)
	EDITS := PARSER.Parse()
	this.editOutcomes(EDITS)
	if len(EDITS) > 0 {
		this.PresentIncompleteActions(contexts)
	}
}

func filterContexts(_projection projections.IncompleteActionsByContext, _filter []string) (filtered_ []*projections.Context) {
	if len(_filter) == 0 {
		return _projection.Contexts
	}
	for _, CONTEXT := range _projection.Contexts {
		log.Println(CONTEXT.Name)
		if CONTEXT.Name == "" {
			filtered_ = append(filtered_, CONTEXT)
		} else {
			for _, context := range _filter {
				if strings.ToLower(CONTEXT.Name) == strings.ToLower(context) {
					filtered_ = append(filtered_, CONTEXT)
					break
				}
			}
		}
	}
	return filtered_
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
	PROJECTION := replayIncompleteActions(this.reader)
	CONTEXTS := filterContexts(PROJECTION, contexts)
	SELECTED := selectNextAction(CONTEXTS)
	RESULT := this.editor.EditTempFile(ux.FormatIncompleteActions(SELECTED))
	PARSER := ux.NewIncompleteActionsParser(this.handler, RESULT, CONTEXTS...)
	EDITS := PARSER.Parse()
	this.editOutcomes(EDITS)
	if len(EDITS) > 0 {
		this.PresentIncompleteAction(contexts)
	}
}

func selectNextAction(contexts_ []*projections.Context) (result_ *projections.Context) {
	result_ = &projections.Context{Name: "Next"}
	var ACTIONS []*projections.ContextualAction
	for _, CONTEXT := range contexts_ {
		for _, ACTION := range CONTEXT.Actions {
			ACTIONS = append(ACTIONS, ACTION)
		}
	}
	if len(ACTIONS) == 0 {
		return result_
	}
	INDEX := rand.Intn(len(ACTIONS))
	result_.Actions = append(result_.Actions, ACTIONS[INDEX])
	return result_
}

func (this *Application) PresentContexts() {
	PROJECTION := replayIncompleteActions(this.reader)
	for _, context := range PROJECTION.Contexts {
		fmt.Printf("- %s (%d)\n", context.Name, len(context.Actions))
	}
}
