package ux

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/ui"
)

func TestOutcomeDetailParserFixture(t *testing.T) {
	gunit.Run(new(OutcomeDetailParserFixture), t)
}

type OutcomeDetailParserFixture struct {
	*gunit.Fixture
	handler   *OutcomeDetailParserFixtureFakeHandler
	outcomeID string
	actionIDs map[string]bool
	content   string
}

func (this *OutcomeDetailParserFixture) Setup() {
	this.handler = NewOutcomeDetailParserFixtureFakeHandler()
	this.actionIDs = make(map[string]bool)
}

func (this *OutcomeDetailParserFixture) parseAndAssertResult(expected ...interface{}) {
	parser := NewOutcomeDetailParser(this.handler, this.outcomeID, this.actionIDs, this.content)

	err := parser.Parse()

	this.So(err, should.BeNil)
	this.So(this.handler.handled, should.Resemble, expected)
}

func (this *OutcomeDetailParserFixture) TestNoChange_NoEvents_NoError() {
	this.content = ui.TrackOutcomeTemplate
	this.parseAndAssertResult()
}

func (this *OutcomeDetailParserFixture) TestTrackNewOutcome_HappyPath() {
	this.outcomeID = ""
	this.content = allElementsAndAllNewTasks

	this.parseAndAssertResult(
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

func (this *OutcomeDetailParserFixture) TestUpdateExistingOutcome() {
	this.content = allElementsAndAllExistingTasks
	this.outcomeID = "0"
	this.actionIDs["1"] = true
	this.actionIDs["2"] = true
	this.actionIDs["3"] = true
	this.actionIDs["4"] = true
	this.actionIDs["5"] = true
	this.actionIDs["6"] = true

	this.parseAndAssertResult(
		&commands.UpdateOutcomeTitle{UpdatedTitle: "The Title"},
		&commands.UpdateOutcomeExplanation{OutcomeID: "0", UpdatedExplanation: "The Explanation"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "1",
			UpdatedDescription: "concurrent complete   @context1 @context2",
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "1"},
		&commands.MarkActionStatusComplete{OutcomeID: "0", ActionID: "1"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "2",
			UpdatedDescription: "concurrent incomplete @context1 @context2",
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "2"},
		&commands.MarkActionStatusIncomplete{OutcomeID: "0", ActionID: "2"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "3",
			UpdatedDescription: "concurrent latent     @context1 @context2",
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "3"},
		&commands.MarkActionStatusLatent{OutcomeID: "0", ActionID: "3"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "4",
			UpdatedDescription: "sequential complete   @context1 @context2",
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "4"},
		&commands.MarkActionStatusComplete{OutcomeID: "0", ActionID: "4"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "5",
			UpdatedDescription: "sequential incomplete @context1 @context2",
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "5"},
		&commands.MarkActionStatusIncomplete{OutcomeID: "0", ActionID: "5"},

		&commands.UpdateActionDescription{
			OutcomeID:          "0",
			ActionID:           "6",
			UpdatedDescription: "sequential latent     @context1 @context2",
		},
		&commands.MarkActionStrategySequential{OutcomeID: "0", ActionID: "6"},
		&commands.MarkActionStatusLatent{OutcomeID: "0", ActionID: "6"},

		&commands.UpdateOutcomeDescription{
			OutcomeID:          "0",
			UpdatedDescription: "\nThe Description\n",
		},
	)
}

func (this *OutcomeDetailParserFixture) TestActionTamperedWith() {
	_, _ = this.handler.Next() // get "0" out of the way
	this.content = tamperedAction
	this.outcomeID = "0"

	this.parseAndAssertResult(
		&commands.UpdateOutcomeTitle{UpdatedTitle: "The Title"},
		&commands.UpdateOutcomeExplanation{OutcomeID: "0", UpdatedExplanation: "The Explanation"},

		&commands.TrackAction{
			OutcomeID:   "0",
			Description: "concurrent complete   @context1 @context2",
			Result:      commands.CreateResult{ID: "1"},
		},
		&commands.MarkActionStrategyConcurrent{OutcomeID: "0", ActionID: "1"},
		&commands.MarkActionStatusIncomplete{OutcomeID: "0", ActionID: "1"},

		&commands.UpdateOutcomeDescription{
			OutcomeID:          "0",
			UpdatedDescription: "\nThe Description\n",
		},
	)
}

func (this *OutcomeDetailParserFixture) TestDeletedAction() {
	this.content = actionDeleted
	this.outcomeID = "0"
	this.actionIDs["1"] = true

	this.parseAndAssertResult(
		&commands.UpdateOutcomeTitle{UpdatedTitle: "The Title"},
		&commands.UpdateOutcomeExplanation{OutcomeID: "0", UpdatedExplanation: "The Explanation"},

		&commands.DeleteAction{OutcomeID: "0", ActionID: "1"},

		&commands.UpdateOutcomeDescription{OutcomeID: "0", UpdatedDescription: "\nThe Description\n"},
	)
}

var (
	allElementsAndAllNewTasks = `# The Title

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

	allElementsAndAllExistingTasks = fmt.Sprintf(`# The Title

> The Explanation


## Actions:

- [X]   %s concurrent complete   @context1 @context2
-  []   %s concurrent incomplete @context1 @context2
-  [?]  %s concurrent latent     @context1 @context2
1. [x]  %s sequential complete   @context1 @context2
2. [ ]  %s sequential incomplete @context1 @context2
10. [?] %s sequential latent     @context1 @context2


## Support Materials:

The Description
`,
		"`0x1`",
		"`0x2`",
		"`0x3`",
		"`0x4`",
		"`0x5`",
		"`0x6`",
	)

	tamperedAction = fmt.Sprintf(`# The Title

> The Explanation


## Actions:

- [ ] %s concurrent complete   @context1 @context2


## Support Materials:

The Description
`,
		"`0xTAMPERED`",
	)

	actionDeleted = `# The Title

> The Explanation


## Actions:



## Support Materials:

The Description
`
)

/////////////////////////////////////////////////////////////////////

type OutcomeDetailParserFixtureFakeHandler struct {
	handled []interface{}

	id   int
	errs []error
}

func NewOutcomeDetailParserFixtureFakeHandler() *OutcomeDetailParserFixtureFakeHandler {
	return &OutcomeDetailParserFixtureFakeHandler{}
}

func (this *OutcomeDetailParserFixtureFakeHandler) Handle(messages ...interface{}) {
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

func (this *OutcomeDetailParserFixtureFakeHandler) Next() (id string, err error) {
	defer func() { this.id++ }()
	if this.id < len(this.errs) {
		err = this.errs[this.id]
	}
	return fmt.Sprint(this.id), err
}
