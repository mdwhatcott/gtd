package projections

import (
	"log"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/events"
)

type OutcomeDetailsProjector struct{ OutcomeDetails }

func NewOutcomeDetailsProjector() *OutcomeDetailsProjector {
	return &OutcomeDetailsProjector{}
}

func (this *OutcomeDetailsProjector) Projection() interface{} {
	return this.OutcomeDetailsProjection()
}
func (this *OutcomeDetailsProjector) OutcomeDetailsProjection() OutcomeDetails {
	return this.OutcomeDetails
}

func (this *OutcomeDetailsProjector) Apply(messages chan interface{}) {
	for MESSAGE := range messages {
		this.apply(MESSAGE)
	}
}

func (this *OutcomeDetailsProjector) apply(MESSAGE interface{}) {
	switch EVENT := MESSAGE.(type) {

	case events.OutcomeTrackedV1:
		this.ID = EVENT.OutcomeID
		this.Title = EVENT.Title

	case events.OutcomeTitleUpdatedV1:
		this.Title = EVENT.UpdatedTitle

	case events.OutcomeDescriptionUpdatedV1:
		this.Description = EVENT.UpdatedDescription

	case events.OutcomeExplanationUpdatedV1:
		this.Explanation = EVENT.UpdatedExplanation

	case events.OutcomeFixedV1:
		this.Status = core.OutcomeStatusFixed

	case events.OutcomeDeferredV1:
		this.Status = core.OutcomeStatusDeferred

	case events.OutcomeUncertainV1:
		this.Status = core.OutcomeStatusUncertain

	case events.OutcomeAbandonedV1:
		this.Status = core.OutcomeStatusAbandoned

	case events.ActionTrackedV1:
		this.Actions = append(this.Actions, &ActionDetails{
			ID:          EVENT.ActionID,
			Description: EVENT.Description,
			Contexts:    EVENT.Contexts,
			Status:      core.ActionStatusIncomplete,
			Strategy:    core.ActionStrategyConcurrent,
		})

	case events.ActionsReorderedV1:
		this.Actions = this.reorderActions(EVENT, EVENT.ReorderedIDs)

	case events.ActionDescriptionUpdatedV1:
		action := this.getAction(EVENT, EVENT.ActionID)
		action.Contexts = EVENT.UpdatedContexts
		action.Description = EVENT.UpdatedDescription

	case events.ActionStatusMarkedLatentV1:
		this.getAction(EVENT, EVENT.ActionID).Status = core.ActionStatusLatent

	case events.ActionStatusMarkedIncompleteV1:
		this.getAction(EVENT, EVENT.ActionID).Status = core.ActionStatusIncomplete

	case events.ActionStatusMarkedCompleteV1:
		this.getAction(EVENT, EVENT.ActionID).Status = core.ActionStatusComplete

	case events.ActionStrategyMarkedSequentialV1:
		this.getAction(EVENT, EVENT.ActionID).Strategy = core.ActionStrategySequential

	case events.ActionStrategyMarkedConcurrentV1:
		this.getAction(EVENT, EVENT.ActionID).Strategy = core.ActionStrategyConcurrent

	case events.ActionDeletedV1:
		this.deleteAction(this.findAction(EVENT.ActionID))

	}
}

func (this *OutcomeDetailsProjector) reorderActions(event interface{}, newOrder []string) (reordered_ []*ActionDetails) {
	for _, ID := range newOrder {
		ACTION := this.getAction(event, ID)
		if len(ACTION.ID) > 0 {
			reordered_ = append(reordered_, ACTION)
		}
	}
	return reordered_
}

func (this *OutcomeDetailsProjector) getAction(event interface{}, id string) *ActionDetails {
	ID := this.findAction(id)
	if ID < 0 {
		log.Printf("Missing action with ID: %s\n\tEvent: %#v", id, event)
		return new(ActionDetails)
	}
	return this.Actions[ID]
}

func (this *OutcomeDetails) findAction(id string) int {
	for i, action := range this.Actions {
		if action.ID == id {
			return i
		}
	}
	return -1
}

func (this *OutcomeDetailsProjector) deleteAction(i int) {
	this.Actions[i] = nil
	this.Actions = append(this.Actions[:i], this.Actions[i+1:]...)
}
