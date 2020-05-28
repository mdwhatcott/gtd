package ux

import (
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
	this.handler.outcomeID = "OutcomeID"
	this.content = happyPath

	this.assertResult(
		&commands.TrackOutcome{Title: "The Title", Result: commands.CreateResult{ID: "OutcomeID"}},
		&commands.UpdateOutcomeExplanation{OutcomeID: "OutcomeID", UpdatedExplanation: "The Explanation"},
	)
}

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

type FakeHandler struct {
	handled []interface{}

	outcomeID string
	err       error
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{}
}

func (this *FakeHandler) Handle(messages ...interface{}) {
	this.handled = append(this.handled, messages...)
	for _, message := range messages {
		switch message := message.(type) {
		case *commands.TrackOutcome:
			message.Result.ID = this.outcomeID
			message.Result.Error = this.err
		}
	}
}

///////////////////////////////////////////////////////////////

const (
	happyPath = `# The Title

> The Explanation


## Actions:

-  [X] concurrent complete   @context1 @context2
-  [ ] concurrent incomplete @context1 @context2
-  [?] concurrent latent     @context1 @context2
1. [X] sequential complete   @context1 @context2
1. [ ] sequential incomplete @context1 @context2
1. [?] sequential latent     @context1 @context2


## Support Materials:

The Description
`
)
