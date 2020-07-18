package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/projections"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/ui"
	"github.com/mdwhatcott/gtd/v3/ui/ux"
)

type Application struct {
	handler core.Handler
	editor  ui.Editor
	reader  joyride.StorageReader
}

func (this *Application) editOutcomes(ids []string) {
	WAITER := new(sync.WaitGroup)
	WAITER.Add(len(ids))
	for _, ID := range ids {
		go this.editOutcome(ID, WAITER)
	}
	WAITER.Wait()
}

func (this *Application) editOutcome(id string, waiter *sync.WaitGroup) {
	defer waiter.Done()

	STREAM := &storage.OutcomeEventStream{OutcomeID: id}
	this.reader.Read(context.Background(), STREAM)
	PROJECTOR := projections.NewOutcomeDetailsProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
	PROJECTION := PROJECTOR.OutcomeDetailsProjection()
	RESULT := this.editor.EditTempFile(ux.FormatOutcomeDetails(PROJECTION))
	PARSER := ux.NewOutcomeDetailParser(this.handler, id, PROJECTION, RESULT)
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
func replayOutcomesListing(reader joyride.StorageReader) projections.OutcomesListing {
	STREAM := &storage.EventStream{}
	reader.Read(context.Background(), STREAM)
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

func filterContexts(projection projections.IncompleteActionsByContext, filter []string) (filtered_ []*projections.Context) {
	if len(filter) == 0 {
		return projection.Contexts
	}
	for _, CONTEXT := range projection.Contexts {
		if CONTEXT.Name == "" {
			filtered_ = append(filtered_, CONTEXT)
		} else {
			for _, context := range filter {
				if strings.ToLower(CONTEXT.Name) == strings.ToLower(context) {
					filtered_ = append(filtered_, CONTEXT)
					break
				}
			}
		}
	}
	return filtered_
}

func replayIncompleteActions(reader joyride.StorageReader) projections.IncompleteActionsByContext {
	STREAM := &storage.EventStream{}
	reader.Read(context.Background(), STREAM)
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

func selectNextAction(contexts []*projections.Context) (result_ *projections.Context) {
	result_ = &projections.Context{Name: "Next"}
	var ACTIONS []*projections.ContextualAction
	for _, CONTEXT := range contexts {
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
	for _, CONTEXT := range PROJECTION.Contexts {
		fmt.Printf("- %s (%d)\n", CONTEXT.Name, len(CONTEXT.Actions))
	}
}
