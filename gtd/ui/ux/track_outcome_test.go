package ux

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

func TestTrackOutcomeFixture(t *testing.T) {
	gunit.Run(new(TrackOutcomeFixture), t)
}

type TrackOutcomeFixture struct {
	*gunit.Fixture
	handler *FakeHandler
	content string
}

func (this *TrackOutcomeFixture) Setup() {
	this.handler = NewFakeHandler()
}

func (this *TrackOutcomeFixture) assertResult(expected ...interface{}) {
	editor := NewFakeEditor()
	editor.resultContent = this.content
	ux := NewTrackOutcomeExperience(this.handler, editor)

	err := ux.Engage()

	this.So(err, should.BeNil)
	this.So(this.handler.handled, should.Resemble, expected)
}

func (this *TrackOutcomeFixture) TestNoChange_NoEvents_NoError() {
	this.content = trackOutcomeTemplate
	this.assertResult()
}

func (this *TrackOutcomeFixture) TestHappyPath() {
	this.content = happyPath

	this.assertResult(
		&commands.TrackOutcome{Title: "The Title", Result: commands.CreateResult{ID: "0"}},
		&commands.UpdateOutcomeExplanation{OutcomeID: "0", UpdatedExplanation: "The Explanation"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "concurrent complete   @context1 @context2",
			Result:      commands.CreateResult{ID: "1"},
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "1"},
		&commands.MarkActionStatusComplete{OutcomeID: "0", ActionID: "1"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "concurrent incomplete @context1 @context2",
			Result:      commands.CreateResult{ID: "2"},
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "2"},
		&commands.MarkActionStatusIncomplete{OutcomeID: "0", ActionID: "2"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "concurrent latent     @context1 @context2",
			Result:      commands.CreateResult{ID: "3"},
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "3"},
		&commands.MarkActionStatusLatent{OutcomeID: "0", ActionID: "3"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "sequential complete   @context1 @context2",
			Result:      commands.CreateResult{ID: "4"},
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "4"},
		&commands.MarkActionStatusComplete{OutcomeID: "0", ActionID: "4"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "sequential incomplete @context1 @context2",
			Result:      commands.CreateResult{ID: "5"},
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "5"},
		&commands.MarkActionStatusIncomplete{OutcomeID: "0", ActionID: "5"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "sequential latent     @context1 @context2",
			Result:      commands.CreateResult{ID: "6"},
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "6"},
		&commands.MarkActionStatusLatent{OutcomeID: "0", ActionID: "6"},

		&commands.UpdateOutcomeDescription{
			OutcomeID:          "0",
			UpdatedDescription: "\nThe Description\n",
		},
	)
}

const (
	happyPath = `# The Title

> The Explanation


## Actions:

- [X]   concurrent complete   @context1 @context2
-  []   concurrent incomplete @context1 @context2
-  [?]  concurrent latent     @context1 @context2
1. [x]  sequential complete   @context1 @context2
2. [ ]  sequential incomplete @context1 @context2
10. [?] sequential latent     @context1 @context2


## Support Materials:

The Description
`
)

////////////////////////////////////////////////////////////////////

type FakeEditor struct {
	initialContent string
	resultContent  string
}

func NewFakeEditor() *FakeEditor {
	return &FakeEditor{}
}

func (this *FakeEditor) EditTempFile(initialContent string) (resultContent string) {
	this.initialContent = initialContent
	return this.resultContent
}

/////////////////////////////////////////////////////////////////////

type FakeHandler struct {
	handled []interface{}

	id   int
	errs []error
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{}
}

func (this *FakeHandler) Handle(messages ...interface{}) {
	this.handled = append(this.handled, messages...)
	for _, message := range messages {
		switch message := message.(type) {
		case *commands.TrackOutcome:
			message.Result.ID, message.Result.Error = this.Next()
		case *commands.TrackAction:
			message.Result.ID, message.Result.Error = this.Next()
		}
	}
}

func (this *FakeHandler) Next() (id string, err error) {
	defer func() { this.id++ }()
	if this.id < len(this.errs) {
		err = this.errs[this.id]
	}
	return fmt.Sprint(this.id), err
}
