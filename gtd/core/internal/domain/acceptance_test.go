package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/joyride/v2"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture

	id      int
	now     time.Time
	shell   *FakeShell
	handler *Handler
	task    *Task
}

func (this *Fixture) Setup() {
	this.now = time.Now()
	this.shell = NewFakeShell(this)
	this.task = NewTask(this.generateID)
	this.task.clock = clock.Freeze(this.now)
	this.task.log = logging.Capture(this)
}
func (this *Fixture) Teardown() {
	this.assertTransferalOfResultOwnership()
}
func (this *Fixture) assertTransferalOfResultOwnership() {
	alreadyPublished := len(this.task.PendingWrites())
	this.task.publishResults()
	doublyPublished := len(this.task.PendingWrites()) - alreadyPublished
	this.So(doublyPublished, should.Equal, 0)
}
func (this *Fixture) handle(commands ...interface{}) {
	runner := joyride.NewRunner(
		joyride.WithStorageReader(this.shell),
		joyride.WithStorageWriter(this.shell),
	)
	this.handler = NewHandler(runner, this.task)
	this.handler.Handle(commands...)
}
func (this *Fixture) generateID() string {
	this.id++
	return fmt.Sprint(this.id)
}
func (this *Fixture) TestUnrecognizedMessageTypes_JoyrideHandlerPanics() {
	this.So(func() { this.handle(42) }, should.PanicWith, joyride.ErrUnknownType)
	this.So(func() { this.handle(true) }, should.PanicWith, joyride.ErrUnknownType)
}
func (this *Fixture) TestTrackOutcome_PublishOutcomeTracked_ReturnOutcomeID() {
	command := &commands.TrackOutcome{Title: "title"}

	this.handle(command)

	this.So(command.Result, should.Resemble, commands.CreateResult{
		ID:    "1",
		Error: nil,
	})
	this.shell.AssertOutput(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: "1",
			Title:     "title",
		},
	)
}
