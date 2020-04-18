package domain

import (
	"github.com/smartystreets/joyride/v2"
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
	switch message.(type) {
	default:
		return false
	}
	return true
}
