package domain

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/events"
	"github.com/mdwhatcott/gtd/v3/util/fake"
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
	log     *bytes.Buffer
	logger  *log.Logger
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

func (this *Fixture) Setup() {
	this.now = time.Now()
	this.log = new(bytes.Buffer)
	this.logger = log.New(this.log, "", 0)
	this.task = NewTask(this.logger, this.Now, this.generateID)
	this.Joyride = fake.NewJoyride(this.logger)
	this.runner = joyride.NewRunner(
		joyride.WithStorageReader(this.Joyride),
		joyride.WithStorageWriter(this.Joyride),
	)
}
func (this *Fixture) Now() time.Time {
	return this.now
}
func (this *Fixture) handle(command interface{}) {
	this.handler = NewHandler(this.logger, this.Now, this.generateID, this.runner)
	this.handler.Handle(context.Background(), command)
}
func (this *Fixture) generateID() string {
	this.id++
	return fmt.Sprint(this.id)
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
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)
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
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)
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

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
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
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
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
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
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
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeAbandoned{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeDeferred_PublishedOutcomeDeferred() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
	)
	COMMAND := &commands.DeclareOutcomeDeferred{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeDeferredV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeclareOutcomeDeferred_AlreadyDeferred_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeferredV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeDeferred{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeDeferred_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeDeferred{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeUncertain_PublishedOutcomeUncertain() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
	)
	COMMAND := &commands.DeclareOutcomeUncertain{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.OutcomeUncertainV1{
			Timestamp: this.now,
			OutcomeID: "1",
		},
	)
}
func (this *Fixture) TestDeclareOutcomeUncertain_AlreadyUncertain_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeUncertainV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeUncertain{OutcomeID: "1"}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeclareOutcomeUncertain_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.DeclareOutcomeUncertain{OutcomeID: "1"}
	this.handle(COMMAND)
	this.So(COMMAND.Result.Error, should.Equal, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestTrackAction_PublishActionTracked() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.TrackAction{
		OutcomeID:   "outcome",
		Description: "description @context1 @context2 @context1",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result, should.Resemble, commands.CreateResult{
		ID:    "1",
		Error: nil,
	})
	this.AssertOutput(
		events.ActionTrackedV1{
			Timestamp:   this.now,
			OutcomeID:   "outcome",
			ActionID:    "1",
			Description: "description @context1 @context2 @context1",
			Contexts:    []string{"context1", "context2"},
		},
	)
}
func (this *Fixture) TestTrackAction_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.TrackAction{
		OutcomeID:   "1",
		Description: "description",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result, should.Resemble, commands.CreateResult{
		Error: core.ErrOutcomeNotFound,
	})
	this.AssertNoOutput()
}
func (this *Fixture) TestTrackAction_IncrementSequence_PublishActionTracked() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
	)

	COMMAND := &commands.TrackAction{
		OutcomeID:   "outcome",
		Description: "description1",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result, should.Resemble, commands.CreateResult{
		ID:    "1",
		Error: nil,
	})
	this.AssertOutput(
		events.ActionTrackedV1{
			Timestamp:   this.now,
			OutcomeID:   "outcome",
			ActionID:    "1",
			Description: "description1",
			Contexts:    nil,
		},
	)
}
func (this *Fixture) TestUpdateActionDescription_PublishActionDescriptionUpdated() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
	)

	COMMAND := &commands.UpdateActionDescription{
		OutcomeID:          "outcome",
		ActionID:           "action",
		UpdatedDescription: "description @context1",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionDescriptionUpdatedV1{
			Timestamp:          this.now,
			OutcomeID:          "outcome",
			ActionID:           "action",
			UpdatedDescription: "description @context1",
			UpdatedContexts:    []string{"context1"},
		},
	)
}
func (this *Fixture) TestUpdateActionDescription_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.UpdateActionDescription{
		OutcomeID:          "outcome",
		ActionID:           "action",
		UpdatedDescription: "description @context1",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateActionDescription_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.UpdateActionDescription{
		OutcomeID:          "1",
		ActionID:           "action",
		UpdatedDescription: "description @context1",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestUpdateActionDescription_DescriptionUnchanged_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionDescriptionUpdatedV1{
			OutcomeID:          "outcome",
			ActionID:           "action",
			UpdatedDescription: "updated description",
		},
	)

	COMMAND := &commands.UpdateActionDescription{
		OutcomeID:          "outcome",
		ActionID:           "action",
		UpdatedDescription: "updated description",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestReorderActions_PublishActionReordered() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action0",
			Description: "description0",
			Contexts:    nil,
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action1",
			Description: "description1",
			Contexts:    nil,
		},
	)

	COMMAND := &commands.ReorderActions{
		OutcomeID:    "outcome",
		ReorderedIDs: []string{"action1", "action0"},
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionsReorderedV1{
			Timestamp:    this.now,
			OutcomeID:    "outcome",
			ReorderedIDs: []string{"action1", "action0"},
		},
	)
}
func (this *Fixture) TestReorderActions_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.ReorderActions{
		OutcomeID:    "1",
		ReorderedIDs: []string{"action1", "action0"},
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Wrap, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestReorderActions_AnyActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action0",
			Description: "description0",
			Contexts:    nil,
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action1",
			Description: "description1",
			Contexts:    nil,
		},
	)

	COMMAND := &commands.ReorderActions{
		OutcomeID:    "outcome",
		ReorderedIDs: []string{"action1", "action-not-found"},
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Wrap, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestReorderActions_AnyActionMissing_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action0",
			Description: "description0",
			Contexts:    nil,
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action1",
			Description: "description1",
			Contexts:    nil,
		},
	)

	COMMAND := &commands.ReorderActions{
		OutcomeID:    "outcome",
		ReorderedIDs: []string{"action1"},
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Wrap, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestReorderActions_NoActions_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.ReorderActions{
		OutcomeID:    "outcome",
		ReorderedIDs: nil,
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Wrap, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusLatent_PublishActionStatusMarkedLatent() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.MarkActionStatusLatent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionStatusMarkedLatentV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestMarkActionStatusLatent_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.MarkActionStatusLatent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusLatent_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.MarkActionStatusLatent{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusLatent_AlreadyMarkedLatent_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionStatusMarkedLatentV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.MarkActionStatusLatent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusIncomplete_PublishActionStatusMarkedIncomplete() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.MarkActionStatusIncomplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionStatusMarkedIncompleteV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestMarkActionStatusIncomplete_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.MarkActionStatusIncomplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusIncomplete_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.MarkActionStatusIncomplete{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusIncomplete_AlreadyMarkedIncomplete_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionStatusMarkedIncompleteV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.MarkActionStatusIncomplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusComplete_PublishActionStatusMarkedComplete() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.MarkActionStatusComplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionStatusMarkedCompleteV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestMarkActionStatusComplete_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.MarkActionStatusComplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusComplete_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.MarkActionStatusComplete{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStatusComplete_AlreadyMarkedComplete_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionStatusMarkedCompleteV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.MarkActionStatusComplete{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategySequential_PublishActionStrategyMarkedSequential() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.MarkActionStrategySequential{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionStrategyMarkedSequentialV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestMarkActionStrategySequential_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.MarkActionStrategySequential{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategySequential_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.MarkActionStrategySequential{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategySequential_AlreadyMarkedSequential_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionStrategyMarkedSequentialV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.MarkActionStrategySequential{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategyConcurrent_PublishActionStrategyMarkedConcurrent() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.MarkActionStrategyConcurrent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionStrategyMarkedConcurrentV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestMarkActionStrategyConcurrent_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.MarkActionStrategyConcurrent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategyConcurrent_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1",
		events.OutcomeTrackedV1{
			OutcomeID: "1",
			Title:     "title",
		},
		events.OutcomeDeletedV1{
			OutcomeID: "1",
		},
	)

	COMMAND := &commands.MarkActionStrategyConcurrent{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestMarkActionStrategyConcurrent_AlreadyMarkedConcurrent_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionStrategyMarkedConcurrentV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.MarkActionStrategyConcurrent{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeUnchanged)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeleteAction_PublishActionDeleted() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
		},
	)

	COMMAND := &commands.DeleteAction{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.BeNil)
	this.AssertOutput(
		events.ActionDeletedV1{
			Timestamp: this.now,
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)
}
func (this *Fixture) TestDeleteAction_ActionNotFound_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
	)

	COMMAND := &commands.DeleteAction{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeleteAction_OutcomeNotFound_ErrorReturned() {
	this.PrepareReadResults("1")

	COMMAND := &commands.DeleteAction{
		OutcomeID: "1",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrOutcomeNotFound)
	this.AssertNoOutput()
}
func (this *Fixture) TestDeleteAction_AlreadyDeleted_ErrorReturned() {
	this.PrepareReadResults("outcome",
		events.OutcomeTrackedV1{
			OutcomeID: "outcome",
			Title:     "title",
		},
		events.ActionTrackedV1{
			OutcomeID:   "outcome",
			ActionID:    "action",
			Description: "description",
			Contexts:    nil,
		},
		events.ActionDeletedV1{
			OutcomeID: "outcome",
			ActionID:  "action",
		},
	)

	COMMAND := &commands.DeleteAction{
		OutcomeID: "outcome",
		ActionID:  "action",
	}
	this.handle(COMMAND)

	this.So(COMMAND.Result.Error, should.Resemble, core.ErrActionNotFound)
	this.AssertNoOutput()
}
