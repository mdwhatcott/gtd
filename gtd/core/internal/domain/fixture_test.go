package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/joyride/v2"

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
}

func (this *Fixture) Setup() {
	this.now = time.Now()
	this.shell = NewFakeShell(this.Fixture)
}

func (this *Fixture) handle(commands ...interface{}) {
	this.handler = NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(this.shell),
			joyride.WithStorageWriter(this.shell),
			joyride.WithMessageDispatcher(this.shell),
		),
		this.generateID,
	)
	this.handler.clock = clock.Freeze(this.now)
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
	this.So(func() { this.handle(&commands.AddContextToAction{})}, should.NotPanic)
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
	this.So(func() { this.handle(&commands.RedefinedAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.AddContextToAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.RemoveContextFromAction{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.MarkActionComplete{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.MarkActionNotComplete{}) }, should.NotPanic)
	this.So(func() { this.handle(&commands.DeleteAction{}) }, should.NotPanic)
}

func (this *Fixture) TestTrackOutcome() {
	command := &commands.TrackOutcome{
		UserID:     "UserID",
		Definition: "The inertial dampers are fixed",
	}

	this.handle(command)

	this.So(command.Result.OutcomeID, should.Equal, "1")
	this.shell.AssertOutput(events.OutcomeFixedV1{
		Timestamp: this.now,
		UserID:    "UserID",
		OutcomeID: "1",
	})
}
