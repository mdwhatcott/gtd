package domain

import (
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

type Handler struct {
	*joyride.Handler

	task *Task
}

func NewHandler(runner joyride.Runner, task *Task) *Handler {
	this := &Handler{task: task}
	this.Handler = joyride.NewHandler(this, runner)
	this.Handler.Add(this.task)
	return this
}
func (this *Handler) HandleMessage(message interface{}) bool {
	switch message := message.(type) {
	case *commands.ProvideOutcomeExplanation:
		this.task.DefineOutcome(message)
	case *commands.RedefineOutcome:
		this.task.RedefineOutcome(message)
	default:
		return false
	}
	return true
}
