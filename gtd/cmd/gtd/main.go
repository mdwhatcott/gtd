package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	flag.Usage = Usage
	flag.Parse()

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	APP := BuildApplication()

	for {
		fmt.Print("\n-> ")
		directive := strings.Fields(ScanLine())
		switch directive[0] {

		case "projects":
			APP.PresentOutcomesListing(directive[1:])

		case "tasks":
			APP.PresentIncompleteActions(directive[1:])

		case "task":
			APP.PresentIncompleteAction(directive[1:])

		case "contexts":
			APP.PresentContexts()

		case "help":
			Usage()

		case "exit":
			return

		default:
			log.Println("Unrecognized directive:", directive)
		}
	}
}

func ScanLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func Usage() {
	flags := log.Flags()
	log.SetFlags(0)
	defer log.SetFlags(flags)

	log.Println(`Usage of gtd:

This application provides a REPL-style interface for selecting directives.
Here are several examples of commands that can be entered:


# To show this usage documentation:
-> help


# To exit the program:
-> exit


# To present a listing of all active contexts:
-> contexts


# [UNDER CONSTRUCTION] To present a single, random, pending task from an active project:
-> task

## [UNDER CONSTRUCTION] Optionally, draw from tasks matching the provided contexts, home and work:
-> task home work

## [UNDER CONSTRUCTION] Optionally, draw from tasks that have no context:
-> task -


# To present all pending tasks from active projects:
-> tasks

## [UNDER CONSTRUCTION] Optionally, show only those tasks that match the provided contexts, home and work:
-> tasks home work

## [UNDER CONSTRUCTION] Optionally, include tasks that have no specified contexts:
-> tasks home work -


# To present all projects (separated by status):
-> projects

## [UNDER CONSTRUCTION] Optionally, only present projects matching the provided statuses, fixed and deferred:
-> projects fixed deferred

## [UNDER CONSTRUCTION] Optionally, only present projects that have no pending tasks:
-> projects fixed deferred ?


The End`)
}

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
