package domain

import (
	"context"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
)

type Handler struct {
	*joyride.Handler
	log    core.Logger
	clock  core.Clock
	nextID core.IDFunc
}

func NewHandler(
	log core.Logger,
	clock core.Clock,
	nextID core.IDFunc,
	runner joyride.Runner,
) *Handler {
	THIS := &Handler{
		log:    log,
		clock:  clock,
		nextID: nextID,
	}
	THIS.Handler = joyride.NewHandler(THIS, runner)
	return THIS
}

func (this *Handler) HandleMessage(_ context.Context, message interface{}) bool {
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
	TASK := NewTask(this.log, this.clock, this.nextID)
	this.Handler.Add(TASK)
	return TASK
}
