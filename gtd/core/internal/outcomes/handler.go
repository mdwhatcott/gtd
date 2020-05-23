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
	TRACK, OK := _message.(*commands.TrackOutcome)
	if OK {
		this.task.PrepareToTrackOutcome(TRACK)
		return true
	}
	ID, OK := _message.(commands.Identifiable)
	if OK {
		this.task.PrepareInstruction(ID)
		return true
	}
	return false
}
