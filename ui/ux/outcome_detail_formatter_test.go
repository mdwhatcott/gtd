package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func TestOutcomeDetailFormatterFixture(t *testing.T) {
	gunit.Run(new(OutcomeDetailFormatterFixture), t)
}

type OutcomeDetailFormatterFixture struct {
	*gunit.Fixture
	outcome projections.OutcomeDetails
}

func (this *OutcomeDetailFormatterFixture) Test() {
	OUTCOME := projections.OutcomeDetails{
		Title:       "Title",
		Explanation: "Explanation",
		Description: "Description",
		Actions: append(this.outcome.Actions,
			&projections.ActionDetails{
				ID:          "000111",
				Description: "Action 1",
				Status:      core.ActionStatusIncomplete,
				Strategy:    core.ActionStrategyConcurrent,
			},
			&projections.ActionDetails{
				ID:          "000222",
				Description: "Action 2",
				Status:      core.ActionStatusLatent,
				Strategy:    core.ActionStrategyConcurrent,
			},
			&projections.ActionDetails{
				ID:          "000333",
				Description: "Action 3",
				Status:      core.ActionStatusComplete,
				Strategy:    core.ActionStrategyConcurrent,
			},
			&projections.ActionDetails{
				ID:          "000444",
				Description: "Action 4",
				Status:      core.ActionStatusIncomplete,
				Strategy:    core.ActionStrategySequential,
			},
			&projections.ActionDetails{
				ID:          "000555",
				Description: "Action 5",
				Status:      core.ActionStatusLatent,
				Strategy:    core.ActionStrategySequential,
			},
			&projections.ActionDetails{
				ID:          "000666",
				Description: "Action 6",
				Status:      core.ActionStatusComplete,
				Strategy:    core.ActionStrategySequential,
			},
		),
	}

	RESULT := FormatOutcomeDetails(OUTCOME)

	this.So(RESULT, should.Equal, strings.Join([]string{
		"# Title",
		"",
		"> Explanation",
		"",
		"",
		"## Actions:",
		"",
		"-  [ ] `0x0001` Action 1",
		"-  [?] `0x0002` Action 2",
		"-  [X] `0x0003` Action 3",
		"1. [ ] `0x0004` Action 4",
		"1. [?] `0x0005` Action 5",
		"1. [X] `0x0006` Action 6",
		"",
		"",
		"## Support Materials:",
		"",
		"Description",
	}, "\n"))
}
