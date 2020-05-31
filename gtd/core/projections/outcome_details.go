package projections

import (
	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
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

func (this *OutcomeDetailsProjector) Apply(_messages ...interface{}) {
	for _, MESSAGE := range _messages {
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
			this.Actions = this.reorderActions(EVENT.ReorderedIDs)

		case events.ActionDescriptionUpdatedV1:
			action := this.getAction(EVENT.ActionID)
			action.Contexts = EVENT.UpdatedContexts
			action.Description = EVENT.UpdatedDescription

		case events.ActionStatusMarkedLatentV1:
			this.getAction(EVENT.ActionID).Status = core.ActionStatusLatent

		case events.ActionStatusMarkedIncompleteV1:
			this.getAction(EVENT.ActionID).Status = core.ActionStatusIncomplete

		case events.ActionStatusMarkedCompleteV1:
			this.getAction(EVENT.ActionID).Status = core.ActionStatusComplete

		case events.ActionStrategyMarkedSequentialV1:
			this.getAction(EVENT.ActionID).Strategy = core.ActionStrategySequential

		case events.ActionStrategyMarkedConcurrentV1:
			this.getAction(EVENT.ActionID).Strategy = core.ActionStrategyConcurrent

		case events.ActionDeletedV1:
			this.deleteAction(this.findAction(EVENT.ActionID))

		}
	}
}

func (this *OutcomeDetailsProjector) reorderActions(_newOrder []string) (reordered_ []*ActionDetails) {
	for _, ID := range _newOrder {
		reordered_ = append(reordered_, this.getAction(ID))
	}
	return reordered_
}

func (this *OutcomeDetailsProjector) getAction(_id string) *ActionDetails {
	return this.Actions[this.findAction(_id)]
}

func (this *OutcomeDetails) findAction(_id string) int {
	for i, action := range this.Actions {
		if action.ID == _id {
			return i
		}
	}
	panic("SHOULD NOT HAPPEN")
}

func (this *OutcomeDetailsProjector) deleteAction(i int) {
	this.Actions[i] = nil
	this.Actions = append(this.Actions[:i], this.Actions[i+1:]...)
}
