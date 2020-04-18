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

	"github.com/mdwhatcott/gtd/gtd/core"
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
func (this *Fixture) TestDefineOutcome_PublishOutcomeDefined() {
	command := &commands.DefineOutcome{
		Definition: "The inertial dampers are fixed",
	}

	this.handle(command)

	this.So(command.Result.OutcomeID, should.Equal, "1")
	this.shell.AssertOutput(
		events.OutcomeDefinedV1{
			Timestamp:  this.now,
			OutcomeID:  "1",
			Definition: "The inertial dampers are fixed",
		},
	)
}
func (this *Fixture) TestRedefineOutcome_PublishOutcomeRedefined() {
	this.shell.PrepareReadResults("1",
		events.OutcomeDefinedV1{
			OutcomeID:  "1",
			Definition: "old",
		},
	)
	command := &commands.RedefineOutcome{
		OutcomeID:     "1",
		NewDefinition: "new",
	}

	this.handle(command)

	this.So(command.Result.Error, should.BeNil)
	this.shell.AssertOutput(
		events.OutcomeRedefinedV1{
			Timestamp:     this.now,
			OutcomeID:     "1",
			NewDefinition: "new",
		},
	)
}
func (this *Fixture) TestRedefineOutcome_UnrecognizedOutcomeID_ReturnErrorOutcomeNotFound() {
	command := &commands.RedefineOutcome{
		OutcomeID:     "not-there",
		NewDefinition: "new-definition",
	}

	this.handle(command)

	this.So(command.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.shell.AssertNoOutput()
}
func (this *Fixture) TestRedefineOutcome_NoChangeToDefinition_ReturnErrOutcomeUnchanged() {
	this.shell.PrepareReadResults("1",
		events.OutcomeDefinedV1{
			OutcomeID:  "1",
			Definition: "old",
		},
	)
	command := &commands.RedefineOutcome{
		OutcomeID:     "1",
		NewDefinition: "old",
	}

	this.handle(command)

	this.So(command.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.shell.AssertNoOutput()
}
