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

func (this *IncompleteActionsByContextProjector) getOutcome(MESSAGE interface{}) *OutcomeDetailsProjector {
	id := MESSAGE.(events.Identifiable).ID()
	outcome := this.all[id]
	if outcome == nil {
		outcome = NewOutcomeDetailsProjector()
		this.all[id] = outcome
	}
	return outcome
}

func (this *IncompleteActionsByContextProjector) buildProjection() {
	contexts := this.filterOutcomeActionsByContext()
	for _, name := range this.sortedContextNames(contexts) {
		this.Contexts = append(this.Contexts, &Context{
			Name:    name,
			Actions: contexts[name],
		})
	}
}

func (this *IncompleteActionsByContextProjector) filterOutcomeActionsByContext() map[string][]*ContextualAction {
	contexts := make(map[string][]*ContextualAction)

	for _, outcome := range this.all {
		if outcome.Status != core.OutcomeStatusFixed {
			continue
		}

		var firstSequential bool
		for _, action := range outcome.Actions {
			if action.Status != core.ActionStatusIncomplete {
				continue
			}

			if action.Strategy == core.ActionStrategySequential && firstSequential {
				continue
			}

			firstSequential = firstSequential || action.Strategy == core.ActionStrategySequential
			if len(action.Contexts) == 0 {
				contexts[""] = append(contexts[""], &ContextualAction{
					ActionDetails: action,
					OutcomeID:     outcome.ID,
					OutcomeTitle:  outcome.Title,
				})
			}
			for _, context := range action.Contexts {
				contexts[context] = append(contexts[context], &ContextualAction{
					ActionDetails: action,
					OutcomeID:     outcome.ID,
					OutcomeTitle:  outcome.Title,
				})
			}
		}
	}

	this.sortActionsWithinContext(contexts)

	return contexts
}

func (this *IncompleteActionsByContextProjector) sortActionsWithinContext(contexts map[string][]*ContextualAction) {
	for _, context := range contexts {
		sort.Slice(context, func(i, j int) bool {
			return context[i].Description < context[j].Description
		})
	}
}

func (this *IncompleteActionsByContextProjector) sortedContextNames(contexts map[string][]*ContextualAction) (names []string) {
	for name := range contexts {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}
