package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func TestIncompleteActionsParserFixture(t *testing.T) {
	gunit.Run(new(IncompleteActionsParserFixture), t)
}

type IncompleteActionsParserFixture struct {
	*gunit.Fixture
	handler    *IncompleteActionsParserFakeHandler
	content    string
	projection projections.IncompleteActionsByContext
}

func (this *IncompleteActionsParserFixture) Setup() {
	this.handler = NewIncompleteActionsParserFakeHandler()
}

func (this *IncompleteActionsParserFixture) TestBetween() {
	this.So(between("==== `0x1234` ====", "`0x", "`"), should.Equal, "1234")
	this.So(between("==== `0x1234` ====", "not-there", "`"), should.Equal, "")
	this.So(between("==== `0x1234` ====", "`0x", "not-there"), should.Equal, "")
}

func (this *IncompleteActionsParserFixture) TestParse() {
	this.content = modifiedIncompleteActionsContent
	action1 := &projections.ActionDetails{ID: "000111", Description: "description 1"}
	action2 := &projections.ActionDetails{ID: "000222", Description: "description 2"}
	action3 := &projections.ActionDetails{ID: "000333", Description: "description 3"}
	action4 := &projections.ActionDetails{ID: "000444", Description: "description 4"}
	action5 := &projections.ActionDetails{ID: "000555", Description: "description 5"}
	action6 := &projections.ActionDetails{ID: "000666", Description: "description 6"}

	this.projection = projections.IncompleteActionsByContext{
		Contexts: []*projections.Context{
			{
				Name: "home",
				Actions: []*projections.ContextualAction{
					{ActionDetails: action1, OutcomeID: "000AAA", OutcomeTitle: "Outcome A"},
					{ActionDetails: action3, OutcomeID: "000BBB", OutcomeTitle: "Outcome B"},
					{ActionDetails: action5, OutcomeID: "000CCC", OutcomeTitle: "Outcome C"},
					{ActionDetails: action6, OutcomeID: "000DDD", OutcomeTitle: "Outcome D"},
				},
			},
			{
				Name: "work",
				Actions: []*projections.ContextualAction{
					{ActionDetails: action2, OutcomeID: "000EEE", OutcomeTitle: "Outcome E"},
					{ActionDetails: action4, OutcomeID: "000FFF", OutcomeTitle: "Outcome F"},
				},
			},
		},
	}

	edits := NewIncompleteActionsParser(this.handler, this.content, this.projection.Contexts...).Parse()

	this.So(edits, should.Resemble, []string{"000BBB", "000FFF"})
	this.So(this.handler.handled, should.Resemble, []interface{}{
		&commands.MarkActionStatusComplete{OutcomeID: "000BBB", ActionID: "000333"},
		&commands.MarkActionStatusLatent{OutcomeID: "000DDD", ActionID: "000666"},
		&commands.MarkActionStatusComplete{OutcomeID: "000EEE", ActionID: "000222"},
	})
}

var modifiedIncompleteActionsContent = strings.Join([]string{
	"## @Home:",
	"",
	"- [ ] `0x0001` description 1 (`0x000A` Outcome A)",
	"\t- [X] `0x0003` description 3 (`0x000B` Outcome B)",
	"- [ ] `0x0005` description 5 (`0x000C` Outcome C)",
	"- [?] `0x0006` description 6 (`0x000D` Outcome D)",
	"",
	"",
	"## @Work:",
	"",
	"- [X] `0x0002` description 2 (`0x000E` Outcome E)",
	"\t- [ ] `0x0004` description 4 (`0x000F` Outcome F)",
}, "\n")

//////////////////////////////////////////////////////////////

type IncompleteActionsParserFakeHandler struct {
	handled []interface{}
}

func NewIncompleteActionsParserFakeHandler() *IncompleteActionsParserFakeHandler {
	return &IncompleteActionsParserFakeHandler{}
}

func (this *IncompleteActionsParserFakeHandler) Handle(messages ...interface{}) {
	this.handled = append(this.handled, messages...)
}
