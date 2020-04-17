package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

type Handler struct {
	*joyride.Handler

	clock  *clock.Clock
	nextID func() string
}

func NewHandler(runner joyride.Runner, id func() string) *Handler {
	this := &Handler{nextID: id}
	this.Handler = joyride.NewHandler(this, runner)
	return this
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
	case *commands.RedefineAction:
	case *commands.AddContextToAction:
	case *commands.RemoveContextFromAction:
	case *commands.MarkActionComplete:
	case *commands.MarkActionNotComplete:
	case *commands.DeleteAction:
	default:
		return false
	}
	this.Add(NewTask(this.clock.UTCNow(), this.nextID, message))
	return true
}
