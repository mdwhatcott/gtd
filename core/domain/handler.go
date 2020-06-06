package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/commands"
)

type Handler struct {
	*joyride.Handler

	clock  *clock.Clock
	nextID core.IDFunc
}

func NewHandler(_runner joyride.Runner, _nextID core.IDFunc) *Handler {
	THIS := &Handler{nextID: _nextID}
	THIS.Handler = joyride.NewHandler(THIS, _runner)
	return THIS
}

func (this *Handler) HandleMessage(_message interface{}) bool {
	TRACK, OK := _message.(*commands.TrackOutcome)
	if OK {
		this.buildTask().PrepareToTrackOutcome(TRACK)
		return true
	}
	ID, OK := _message.(commands.Identifiable)
	if OK {
		this.buildTask().PrepareInstruction(ID)
		return true
	}
	return false
}

func (this *Handler) buildTask() *Task {
	task := NewTask(this.nextID)
	task.clock = this.clock
	this.Handler.Add(task)
	return task
}
