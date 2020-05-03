package outcomes

import (
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

type Handler struct {
	*joyride.Handler

	task *Task
}

func NewHandler(_runner joyride.Runner, _task *Task) *Handler {
	THIS := &Handler{task: _task}
	THIS.Handler = joyride.NewHandler(THIS, _runner)
	THIS.Handler.Add(THIS.task)
	return THIS
}

func (this *Handler) HandleMessage(_message interface{}) bool {
	switch MESSAGE := _message.(type) {
	case *commands.TrackOutcome:
		this.task.PrepareToTrackOutcome(MESSAGE)
	case *commands.UpdateOutcomeTitle:
		this.task.PrepareInstruction(MESSAGE, MESSAGE.OutcomeID)
	case *commands.UpdateOutcomeExplanation:
		this.task.PrepareInstruction(MESSAGE, MESSAGE.OutcomeID)
	case *commands.UpdateOutcomeDescription:
		this.task.PrepareInstruction(MESSAGE, MESSAGE.OutcomeID)
	default:
		return false
	}
	return true
}
