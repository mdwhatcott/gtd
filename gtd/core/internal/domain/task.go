package domain

import (
	"time"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Task struct {
	*joyride.Base

	now     time.Time
	nextID  func() string
	command interface{}
}

func NewTask(now time.Time, nextID func() string, command interface{}) *Task {
	return &Task{
		Base:    joyride.New(),
		now:     now,
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
	case *commands.RedefineOutcome:
		this.redefineOutcome(command)
	}
}

func (this *Task) trackOutcome(command *commands.TrackOutcome) {
	command.Result.OutcomeID = this.nextID()
	this.AddPendingWrites(
		events.OutcomeDefinedV1{
			Timestamp: this.now,
			OutcomeID: command.Result.OutcomeID,
		},
	)
}

func (this *Task) redefineOutcome(command *commands.RedefineOutcome) {
	this.AddPendingWrites(
		events.OutcomeRedefinedV1{
			Timestamp:     this.now,
			OutcomeID:     command.OutcomeID,
			NewDefinition: command.NewDefinition,
		},
	)
}
