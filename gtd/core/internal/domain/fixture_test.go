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
	gunit.Run(new(Fixture), t, gunit.Options.AllSequential())
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
}
func (this *Fixture) handle(commands ...interface{}) {
	this.task = NewTask(this.generateID)
	this.task.clock = clock.Freeze(this.now)
	this.task.log = logging.Capture(this)
	this.handler = NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(this.shell),
			joyride.WithStorageWriter(this.shell),
		),
		this.task,
	)
	this.handler.Handle(commands...)
}
func (this *Fixture) generateID() string {
	this.id++
	return fmt.Sprint(this.id)
}
func (this *Fixture) TestHandlerPanicsOnUnrecognizedMessageTypes() {
	this.So(func() { this.handle(42) }, should.PanicWith, joyride.ErrUnknownType)
	this.So(func() { this.handle(true) }, should.PanicWith, joyride.ErrUnknownType)
}
func (this *Fixture) TestHandlerAcceptsKnownMessageType() {
	this.So(func() { this.handle(&commands.AddContextToAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.RedefineOutcome{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DescribeOutcome{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeleteOutcome{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeclareOutcomeFixed{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeclareOutcomeRealized{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeclareOutcomeAbandoned{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeclareOutcomeDeferred{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeclareOutcomeUncertain{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.TrackAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.ResequencedAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.RedefineAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.AddContextToAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.RemoveContextFromAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.MarkActionComplete{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.MarkActionNotComplete{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeleteAction{}) }, should.NotPanic)
}
func (this *Fixture) TestDefineOutcome() {
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
func (this *Fixture) TestOutcomeRedefined() {
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
func (this *Fixture) TestRedefineOutcome_NoChangeToDefinition_NoEventPublished() {
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

	this.So(command.Result.Error, should.NotBeNil)
	this.shell.AssertOutput()
}
func (this *Fixture) SkipTestRedefineOutcome_UnrecognizedOutcomeID_ErrorReturnedOnCommand() {

}
