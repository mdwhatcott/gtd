package outcomes

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

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/util/fake"
)

func TestFixture(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture
	*fake.Joyride

	id      int
	now     time.Time
	log     *logging.Logger
	handler *Handler
	task    *Task
	runner  joyride.Runner
}

func (this *Fixture) AssertNoOutput() {
	this.AssertOutput()
}
func (this *Fixture) AssertOutput(_expected ...interface{}) {
	this.So(this.Writes, should.Resemble, _expected)
}
func (this *Fixture) assertTransferalOfResultOwnership() {
	alreadyPUBLISHED := len(this.task.PendingWrites())
	this.task.publishResults()
	doublyPUBLISHED := len(this.task.PendingWrites()) - alreadyPUBLISHED
	this.So(doublyPUBLISHED, should.Equal, 0)
}

func (this *Fixture) Setup() {
	this.now = time.Now()
	this.log = logging.Capture(ioutil.Discard)
	this.log.SetFlags(log.Lshortfile)
	this.log.SetPrefix("--> ")
	this.task = NewTask(this.generateID)
	this.task.clock = clock.Freeze(this.now)
	this.task.log = this.log
	this.Joyride = fake.NewJoyride(this.log)
	this.runner = joyride.NewRunner(
		joyride.WithStorageReader(this.Joyride),
		joyride.WithStorageWriter(this.Joyride),
	)
}
func (this *Fixture) Teardown() {
	this.assertTransferalOfResultOwnership()
}
func (this *Fixture) handle(command interface{}) {
	this.handler = NewHandler(this.runner, this.task)
	this.handler.Handle(command)
}
func (this *Fixture) generateID() string {
	this.id++
	return fmt.Sprint(this.id)
}
func (this *Fixture) enableLogging() {
	this.log.SetOutput(this.Fixture)
}

func (this *Fixture) TestUnrecognizedMessageTypes_JoyrideHandlerPanics() {
	this.So(func() { this.handle(42) }, should.PanicWith, joyride.ErrUnknownType)
	this.So(func() { this.handle(true) }, should.PanicWith, joyride.ErrUnknownType)
}
func (this *Fixture) TestTrackOutcome_PublishOutcomeTrackedAndFixed_ReturnOutcomeID() {
	COMMAND := &commands.TrackOutcome{Title: "title"}

	this.handle(COMMAND)

	this.So(COMMAND.Result, should.Resemble, commands.CreateResult{ID: "1"})
	this.AssertOutput(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeFixedV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeclareOutcomeFixed_AlreadyFixed_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeFixedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeFixed{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeFixed_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeFixedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeFixed{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeFixed_AfterOutcomeRealized_PublishOutcomeFixed() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeRealizedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeFixed{OutcomeID: "1"}
	this.handle(COMMAND)

	this.AssertOutput(
		events.OutcomeFixedV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestUpdateOutcomeTitle_PublishOutcomeTitleUpdated() {
	this.PrepareReadResults("1", events.OutcomeTrackedV1{
		OutcomeID: "1",
		Title:     "title",
	})
	COMMAND := &commands.UpdateOutcomeTitle{
		OutcomeID:    "1",
		UpdatedTitle: "new-title",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeTitleUpdatedV1{
			Timestamp:    this.now,
			OutcomeID:    "1",
			UpdatedTitle: "new-title",
		},
	)
}
func (this *Fixture) TestUpdateOutcomeTitle_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1", nil)
	COMMAND := &commands.UpdateOutcomeTitle{
		OutcomeID:    "1",
		UpdatedTitle: "new-title",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeTitle_ContentUnchangedSinceCreation_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "first-title",
		},
	)
	COMMAND := &commands.UpdateOutcomeTitle{
		OutcomeID:    "1",
		UpdatedTitle: "first-title",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeTitle_ContentUnchangedSinceLastUpdate_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "first-title",
		},
		events.OutcomeTitleUpdatedV1{
			OutcomeID:    "1",
			UpdatedTitle: "second-title",
		},
	)
	COMMAND := &commands.UpdateOutcomeTitle{
		OutcomeID:    "1",
		UpdatedTitle: "second-title",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeExplanation_PublishOutcomeExplanationUpdated() {
	this.PrepareReadResults("1", events.OutcomeTrackedV1{
		OutcomeID: "1",
		Title:     "title",
	})
	COMMAND := &commands.UpdateOutcomeExplanation{
		OutcomeID:          "1",
		UpdatedExplanation: "explanation",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeExplanationUpdatedV1{
			Timestamp:          this.now,
			OutcomeID:          "1",
			UpdatedExplanation: "explanation",
		},
	)
}
func (this *Fixture) TestUpdateOutcomeExplanation_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1", nil)
	COMMAND := &commands.UpdateOutcomeExplanation{
		OutcomeID:          "1",
		UpdatedExplanation: "new-explanation",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeExplanation_ContentUnchanged_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeExplanationUpdatedV1{
			OutcomeID:          "1",
			UpdatedExplanation: "first-explanation",
		},
	)
	COMMAND := &commands.UpdateOutcomeExplanation{
		OutcomeID:          "1",
		UpdatedExplanation: "first-explanation",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeDescription_PublishOutcomeDescriptionUpdated() {
	this.PrepareReadResults("1", events.OutcomeTrackedV1{
		OutcomeID: "1",
		Title:     "title",
	})
	COMMAND := &commands.UpdateOutcomeDescription{
		OutcomeID:          "1",
		UpdatedDescription: "explanation",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeDescriptionUpdatedV1{
			Timestamp:          this.now,
			OutcomeID:          "1",
			UpdatedDescription: "explanation",
		},
	)
}
func (this *Fixture) TestUpdateOutcomeDescription_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1", nil)
	COMMAND := &commands.UpdateOutcomeDescription{
		OutcomeID:          "1",
		UpdatedDescription: "new-description",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateOutcomeDescription_ContentUnchanged_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDescriptionUpdatedV1{
			OutcomeID:          "1",
			UpdatedDescription: "first-description",
		},
	)
	COMMAND := &commands.UpdateOutcomeDescription{
		OutcomeID:          "1",
		UpdatedDescription: "first-description",
	}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeleteOutcome_PublishedOutcomeDeleted() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
	)
	COMMAND := &commands.DeleteOutcome{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeDeletedV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeleteOutcome_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1", nil)
	COMMAND := &commands.DeleteOutcome{OutcomeID: "1"}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeleteOutcome_AlreadyDeleted_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeleteOutcome{OutcomeID: "1"}

	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeRealized_PublishedOutcomeRealized() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
	)
	COMMAND := &commands.DeclareOutcomeRealized{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeRealizedV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeclareOutcomeRealized_AlreadyRealized_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeRealizedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeRealized{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeRealized_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeRealizedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeRealized{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeAbandoned_PublishedOutcomeAbandoned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
	)
	COMMAND := &commands.DeclareOutcomeAbandoned{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeAbandonedV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeclareOutcomeAbandoned_AlreadyAbandoned_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeAbandonedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeAbandoned{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeAbandoned_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeAbandonedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeAbandoned{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
