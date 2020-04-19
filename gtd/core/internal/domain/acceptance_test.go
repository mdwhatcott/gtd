package domain

import (
	"fmt"
	"io/ioutil"
	"log"
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
	*FakeShell

	id      int
	now     time.Time
	log     *logging.Logger
	handler *Handler
	task    *Task
	runner  joyride.Runner
}

func (this *Fixture) Setup() {
	this.now = time.Now()
	this.log = logging.Capture(ioutil.Discard)
	this.log.SetFlags(log.Lshortfile)
	this.log.SetPrefix("--> ")
	this.task = NewTask(this.generateID)
	this.task.clock = clock.Freeze(this.now)
	this.task.log = this.log
	this.FakeShell = NewFakeShell(this)
	this.runner = joyride.NewRunner(
		joyride.WithStorageReader(this.FakeShell),
		joyride.WithStorageWriter(this.FakeShell),
	)
}
func (this *Fixture) Teardown() {
	this.assertTransferalOfResultOwnership()
}
func (this *Fixture) enableLogging() {
	this.log.SetOutput(this.Fixture)
}
func (this *Fixture) assertTransferalOfResultOwnership() {
	alreadyPublished := len(this.task.PendingWrites())
	this.task.publishResults()
	doublyPublished := len(this.task.PendingWrites()) - alreadyPublished
	this.So(doublyPublished, should.Equal, 0)
}
func (this *Fixture) handle(command interface{}) {
	this.handler = NewHandler(this.runner, this.task)
	this.handler.Handle(command)
}
func (this *Fixture) generateID() string {
	this.id++
	return fmt.Sprint(this.id)
}
func (this *Fixture) AssertEventOutput(expected ...interface{}) {
	this.So(this.task.PendingWrites(), should.Resemble, expected)
}
func (this *Fixture) TestUnrecognizedMessageTypes_JoyrideHandlerPanics() {
	this.So(func() { this.handle(42) }, should.PanicWith, joyride.ErrUnknownType)
	this.So(func() { this.handle(true) }, should.PanicWith, joyride.ErrUnknownType)
}
func (this *Fixture) TestTrackOutcome_PublishOutcomeTracked_ReturnOutcomeID() {
	command := &commands.TrackOutcome{Title: "title"}

	this.handle(command)

	this.So(command.Result, should.Resemble, commands.CreateResult{ID: "1"})
	this.AssertOutput(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: "1",
			Title:     "title",
		},
	)
}
func (this *Fixture) TestProvideOutcomeExplanation_PublishOutcomeExplanationProvided() {
	this.PrepareReadResults("1", events.OutcomeTrackedV1{
		OutcomeID: "1",
		Title:     "title",
	})
	command := &commands.ProvideOutcomeExplanation{
		OutcomeID:   "1",
		Explanation: "explanation",
	}

	this.handle(command)

	this.So(command.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeExplanationProvidedV1{
			Timestamp:   this.now,
			OutcomeID:   "1",
			Explanation: "explanation",
		},
	)
}
