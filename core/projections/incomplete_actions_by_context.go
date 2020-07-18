package projections

import (
	"sort"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/events"
)

type IncompleteActionsByContextProjector struct {
	IncompleteActionsByContext
	all map[string]*OutcomeDetailsProjector
}

func NewIncompleteActionsByContextProjector() *IncompleteActionsByContextProjector {
	return &IncompleteActionsByContextProjector{
		all: make(map[string]*OutcomeDetailsProjector),
	}
}

func (this *IncompleteActionsByContextProjector) IncompleteActionsByContextProjection() IncompleteActionsByContext {
	return this.IncompleteActionsByContext
}

func (this *IncompleteActionsByContextProjector) Projection() interface{} {
	return this.IncompleteActionsByContextProjection()
}

func (this *IncompleteActionsByContextProjector) Apply(_messages chan interface{}) {
	defer this.buildProjection()

	for MESSAGE := range _messages {
		this.getOutcome(MESSAGE).apply(MESSAGE)

		switch MESSAGE := MESSAGE.(type) {
		case events.OutcomeDeletedV1:
			delete(this.all, MESSAGE.OutcomeID)
		}
	}
}

func (this *IncompleteActionsByContextProjector) getOutcome(message interface{}) *OutcomeDetailsProjector {
	id := message.(events.Identifiable).ID()
	OUTCOME := this.all[id]
	if OUTCOME == nil {
		OUTCOME = NewOutcomeDetailsProjector()
		this.all[id] = OUTCOME
	}
	return OUTCOME
}

func (this *IncompleteActionsByContextProjector) buildProjection() {
	CONTEXTS := this.filterOutcomeActionsByContext()
	for _, NAME := range this.sortedContextNames(CONTEXTS) {
		this.Contexts = append(this.Contexts, &Context{
			Name:    NAME,
			Actions: CONTEXTS[NAME],
		})
	}
}

func (this *IncompleteActionsByContextProjector) filterOutcomeActionsByContext() (contexts_ map[string][]*ContextualAction) {
	contexts_ = make(map[string][]*ContextualAction)

	for _, OUTCOME := range this.all {
		if OUTCOME.Status != core.OutcomeStatusFixed {
			continue
		}

		var sequential bool
		for _, ACTION := range OUTCOME.Actions {
			if ACTION.Status != core.ActionStatusIncomplete {
				continue
			}

			if ACTION.Strategy == core.ActionStrategySequential && sequential {
				continue
			}

			sequential = sequential || ACTION.Strategy == core.ActionStrategySequential
			if len(ACTION.Contexts) == 0 {
				contexts_[""] = append(contexts_[""], &ContextualAction{
					ActionDetails: ACTION,
					OutcomeID:     OUTCOME.ID,
					OutcomeTitle:  OUTCOME.Title,
				})
			}
			for _, context := range ACTION.Contexts {
				contexts_[context] = append(contexts_[context], &ContextualAction{
					ActionDetails: ACTION,
					OutcomeID:     OUTCOME.ID,
					OutcomeTitle:  OUTCOME.Title,
				})
			}
		}
	}

	this.sortActionsWithinContext(contexts_)

	return contexts_
}

func (this *IncompleteActionsByContextProjector) sortActionsWithinContext(contexts map[string][]*ContextualAction) {
	for _, CONTEXT := range contexts {
		sort.Slice(CONTEXT, func(i, j int) bool {
			return CONTEXT[i].Description < CONTEXT[j].Description
		})
	}
}

func (this *IncompleteActionsByContextProjector) sortedContextNames(contexts map[string][]*ContextualAction) (names_ []string) {
	for NAME := range contexts {
		names_ = append(names_, NAME)
	}
	sort.Slice(names_, func(i, j int) bool {
		return names_[i] < names_[j]
	})
	return names_
}
