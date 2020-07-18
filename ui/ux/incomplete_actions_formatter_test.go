package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func TestIncompleteActionsFormatterFixture(t *testing.T) {
	gunit.Run(new(IncompleteActionsFormatterFixture), t)
}

type IncompleteActionsFormatterFixture struct {
	*gunit.Fixture
}

func (this *IncompleteActionsFormatterFixture) Test() {
	action1 := &projections.ActionDetails{ID: "000111", Description: "description 1"}
	action2 := &projections.ActionDetails{ID: "000222", Description: "description 2"}
	action3 := &projections.ActionDetails{ID: "000333", Description: "description 3"}
	action4 := &projections.ActionDetails{ID: "000444", Description: "description 4"}
	action5 := &projections.ActionDetails{ID: "000555", Description: "description 5"}
	action6 := &projections.ActionDetails{ID: "000666", Description: "description 6"}

	projection := projections.IncompleteActionsByContext{
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

	result := FormatIncompleteActions(projection.Contexts...)

	this.So(result, should.Equal, strings.Join([]string{
		"## @Home:",
		"",
		"- [ ] `0x0001` description 1 (`0x000A` Outcome A)",
		"- [ ] `0x0003` description 3 (`0x000B` Outcome B)",
		"- [ ] `0x0005` description 5 (`0x000C` Outcome C)",
		"- [ ] `0x0006` description 6 (`0x000D` Outcome D)",
		"",
		"",
		"## @Work:",
		"",
		"- [ ] `0x0002` description 2 (`0x000E` Outcome E)",
		"- [ ] `0x0004` description 4 (`0x000F` Outcome F)",
	}, "\n"))
}
