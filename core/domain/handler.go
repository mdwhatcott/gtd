package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
)

type Handler struct {
	*joyride.Handler

	clock  *clock.Clock
	nextID core.IDFunc
}

func NewHandler(runner joyride.Runner, nextID core.IDFunc) *Handler {
	THIS := &Handler{nextID: nextID}
	THIS.Handler = joyride.NewHandler(THIS, runner)
	return THIS
}

func (this *Handler) HandleMessage(message interface{}) bool {
	TRACK, OK := message.(*commands.TrackOutcome)
	if OK {
		this.buildTask().PrepareToTrackOutcome(TRACK)
		return true
	}
	ID, OK := message.(commands.Identifiable)
	if OK {
		this.buildTask().PrepareInstruction(ID)
		return true
	}
	return false
}

func (this *Handler) buildTask() *Task {
	TASK := NewTask(this.nextID)
	TASK.clock = this.clock
	this.Handler.Add(TASK)
	return TASK
}
