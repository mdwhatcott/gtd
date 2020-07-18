package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/projections"
	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/ui"
	"github.com/mdwhatcott/gtd/ui/ux"
)

type Application struct {
	handler core.Handler
	editor  ui.Editor
	reader  joyride.StorageReader

	storageDirectory string
}

func (this *Application) editOutcomes(_ids []string) {
	WAITER := new(sync.WaitGroup)
	WAITER.Add(len(_ids))
	for _, ID := range _ids {
		go this.editOutcome(ID, WAITER)
	}
	WAITER.Wait()
}

func (this *Application) editOutcome(ID string, waiter *sync.WaitGroup) {
	defer waiter.Done()

	STREAM := &storage.OutcomeEventStream{OutcomeID: ID}
	this.reader.Read(STREAM)
	PROJECTOR := projections.NewOutcomeDetailsProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
	PROJECTION := PROJECTOR.OutcomeDetailsProjection()
	RESULT := this.editor.EditTempFile(ux.FormatOutcomeDetails(PROJECTION))
	PARSER := ux.NewOutcomeDetailParser(this.handler, ID, PROJECTION, RESULT)
	ERR := PARSER.Parse()
	if ERR != nil {
		log.Fatal(ERR)
	}
}

func (this *Application) PresentOutcomesListing(statuses []string) {
	PROJECTION := replayOutcomesListing(this.reader)
	PROJECTION = filterOnStatus(PROJECTION, statuses)
	RESULT := this.editor.EditTempFile(ux.FormatOutcomesListing(PROJECTION))
	PARSER := ux.NewOutcomesListingParser(this.handler, PROJECTION, RESULT)
	EDITS := PARSER.Parse()
	this.editOutcomes(EDITS)
	if len(EDITS) > 0 {
		this.PresentOutcomesListing(statuses)
	}
}

func filterOnStatus(projection projections.OutcomesListing, statuses []string) projections.OutcomesListing {
	if len(statuses) == 0 {
		return projection
	}
	if !hasStatus(statuses, "fixed") {
		projection.Fixed = nil
	}
	if !hasStatus(statuses, "deferred") {
		projection.Deferred = nil
	}
	if !hasStatus(statuses, "uncertain") {
		projection.Uncertain = nil
	}
	if !hasStatus(statuses, "abandoned") {
		projection.Abandoned = nil
	}
	return projection
}

func hasStatus(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if strings.ToLower(straw) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}
func replayOutcomesListing(_reader joyride.StorageReader) projections.OutcomesListing {
	STREAM := &storage.EventStream{}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewOutcomesListingProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
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
	STREAM := &storage.EventStream{}
	_reader.Read(STREAM)
	PROJECTOR := projections.NewIncompleteActionsByContextProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
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

func (this *Application) PushChanges() {
	STATUS := exec.Command("git", "status", "--porcelain")
	STATUS.Dir = this.storageDirectory
	OUT, ERR := STATUS.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	if !strings.Contains(strings.TrimSpace(string(OUT)), "events.csv") {
		return
	}

	log.Println("Staging newly generated events...")
	ADD := exec.Command("git", "add", "events.csv")
	ADD.Dir = this.storageDirectory
	OUT, ERR = ADD.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	log.Println("Committing newly generated events...")
	TODAY := time.Now().Format("2006-01-02")
	COMMIT := exec.Command("git", "commit", "-m", TODAY)
	COMMIT.Dir = this.storageDirectory
	OUT, ERR = COMMIT.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	log.Println("Pushing newly generated events...")
	PUSH := exec.Command("git", "push", "origin", "main")
	PUSH.Dir = this.storageDirectory
	OUT, ERR = PUSH.CombinedOutput()
	if ERR != nil {
		log.Println(PUSH)
		log.Fatal(ERR)
	}

	log.Println("Finished.")
}
