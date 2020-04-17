package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Task struct {
	*joyride.Base

	clock   *clock.Clock
	nextID  func() string
	command interface{}
}

func NewTask(clock *clock.Clock, nextID func() string, command interface{}) *Task {
	return &Task{
		Base:    joyride.New(),
		clock:   clock,
		nextID:  nextID,
		command: command,
	}
}

func (this *Task) RequiredReads() []interface{} {
	switch this.command.(type) {
	}
	return this.Base.RequiredReads()
}

func (this *Task) Execute() joyride.TaskResult {
	this.processCommand()
	return this
}

func (this *Task) processCommand() {
	switch command := this.command.(type) {
	case *commands.TrackOutcome:
		this.trackOutcome(command)
	}
}

func (this *Task) trackOutcome(command *commands.TrackOutcome) {
	command.Result.OutcomeID = this.nextID()
	this.AddPendingWrites(events.OutcomeFixedV1{
		Timestamp: this.clock.UTCNow(),
		UserID:    command.UserID,
		OutcomeID: command.Result.OutcomeID,
	})
}
