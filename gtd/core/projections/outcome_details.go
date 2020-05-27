package projections

import "github.com/mdwhatcott/gtd/gtd/core/events"

type OutcomeDetailsProjector struct {
	OutcomeDetails
	//actions map[string]*ActionDetails
}

func NewOutcomeDetailsProjector() *OutcomeDetailsProjector {
	return &OutcomeDetailsProjector{
		//actions: make(map[string]*ActionDetails),
	}
}

func (this *OutcomeDetailsProjector) Projection() interface{} {
	return this.OutcomeDetails
}

func (this *OutcomeDetailsProjector) Apply(_messages ...interface{}) {
	for _, MESSAGE := range _messages {
		switch EVENT := MESSAGE.(type) {

		case events.OutcomeTrackedV1:
			this.Title = EVENT.Title

		case events.OutcomeTitleUpdatedV1:
			this.Title = EVENT.UpdatedTitle

		case events.OutcomeDescriptionUpdatedV1:
			this.Description = EVENT.UpdatedDescription

		case events.OutcomeExplanationUpdatedV1:
			this.Explanation = EVENT.UpdatedExplanation

		case events.ActionTrackedV1:
			action := &ActionDetails{
				ID:          EVENT.ActionID,
				Description: EVENT.Description,
				Contexts:    EVENT.Contexts,
			}
			//this.actions[EVENT.ActionID] = action
			this.Actions = append(this.Actions, action)

		case events.ActionsReorderedV1:

		case events.ActionDescriptionUpdatedV1:
			action := this.Actions[this.findAction(EVENT.ActionID)]
			action.Contexts = EVENT.UpdatedContexts
			action.Description = EVENT.UpdatedDescription

		case events.ActionStatusMarkedLatentV1:

		case events.ActionStatusMarkedIncompleteV1:

		case events.ActionStatusMarkedCompleteV1:

		case events.ActionStrategyMarkedSequentialV1:

		case events.ActionStrategyMarkedConcurrentV1:

		case events.ActionDeletedV1:
			this.deleteAction(this.findAction(EVENT.ActionID))

		}
	}
}

type OutcomeDetails struct {
	Title       string
	Explanation string
	Description string
	Actions     []*ActionDetails
}

func (this *OutcomeDetails) findAction(_id string) int {
	for i, action := range this.Actions {
		if action.ID == _id {
			return i
		}
	}
	panic("SHOULD NOT HAPPEN")
}

func (this *OutcomeDetails) deleteAction(i int) {
	this.Actions[i] = nil
	this.Actions = append(this.Actions[:i], this.Actions[i+1:]...)
}

type ActionDetails struct {
	ID          string
	Description string
	Contexts    []string
}
