package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

type Handler struct {
	*joyride.Handler
	clock *clock.Clock
	tasks map[string]*Task
}

func NewHandler(runner joyride.Runner) *Handler {
	this := &Handler{tasks: make(map[string]*Task)}
	this.Handler = joyride.NewHandler(this, runner)
	return this
}

func (this *Handler) loadTask(userID string) *Task {
	task, found := this.tasks[userID]
	if !found {
		task = NewTask(userID)
		this.tasks[userID] = task
	}
	return task
}

func (this *Handler) HandleMessage(message interface{}) bool {
	switch message.(type) {
	case *commands.TrackOutcome:
	case *commands.RedefineOutcome:
	case *commands.DescribeOutcome:
	case *commands.DeleteOutcome:
	case *commands.DeclareOutcomeFixed:
	case *commands.DeclareOutcomeRealized:
	case *commands.DeclareOutcomeAbandoned:
	case *commands.DeclareOutcomeDeferred:
	case *commands.DeclareOutcomeUncertain:
	case *commands.TrackAction:
	case *commands.ResequencedAction:
	case *commands.RedefinedAction:
	case *commands.AddContextToAction:
	case *commands.RemoveContextFromAction:
	case *commands.MarkActionComplete:
	case *commands.MarkActionNotComplete:
	case *commands.DeleteAction:
	default:
		return false
	}
	return true
}
