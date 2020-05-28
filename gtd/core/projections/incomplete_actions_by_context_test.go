package projections

import (
	"testing"

	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

func TestIncompleteActionsByContextFixture(t *testing.T) {
	gunit.Run(new(IncompleteActionsByContextFixture), t)
}

type IncompleteActionsByContextFixture struct {
	*gunit.Fixture
	*ProjectorFixture
}

func (this *IncompleteActionsByContextFixture) Setup() {
	this.ProjectorFixture = InitializeProjectorFixture(this.Fixture, NewIncompleteActionsByContextProjector())
}

func (this *IncompleteActionsByContextFixture) Test() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "0", Title: "0"},
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "1"},
		events.OutcomeTrackedV1{OutcomeID: "2", Title: "2"},
		events.OutcomeTrackedV1{OutcomeID: "3", Title: "3"},
		events.OutcomeTrackedV1{OutcomeID: "4", Title: "4"},
		events.OutcomeTrackedV1{OutcomeID: "5", Title: "5"},
		events.OutcomeTrackedV1{OutcomeID: "6", Title: "6"},
		events.OutcomeTrackedV1{OutcomeID: "7", Title: "7"},
		events.OutcomeTrackedV1{OutcomeID: "8", Title: "8"},
		events.OutcomeTrackedV1{OutcomeID: "9", Title: "9"},

		events.ActionTrackedV1{
			OutcomeID:   "0",
			ActionID:    "000",
			Description: "action000",
			Contexts:    []string{"c1", "c2"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "0",
			ActionID:    "00",
			Description: "action00",
			Contexts:    []string{"c0", "c1"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "0",
			ActionID:    "0000",
			Description: "action0000",
			Contexts:    []string{"c2", "c3"},
		},

		events.ActionTrackedV1{
			OutcomeID:   "1",
			ActionID:    "11",
			Description: "action11",
			Contexts:    []string{"c0", "c1"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "2",
			ActionID:    "22",
			Description: "action22",
			Contexts:    []string{"c1", "c2"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "3",
			ActionID:    "33",
			Description: "action33",
			Contexts:    []string{"c1", "c2"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "4",
			ActionID:    "44",
			Description: "action44",
			Contexts:    []string{"c1", "c2"},
		},

		events.ActionTrackedV1{
			OutcomeID:   "9",
			ActionID:    "99",
			Description: "action99",
			Contexts:    []string{"c0", "c1"},
		},
		events.ActionTrackedV1{
			OutcomeID:   "9",
			ActionID:    "999",
			Description: "action999",
			Contexts:    []string{"c2", "c3"},
		},

		events.ActionStatusMarkedCompleteV1{OutcomeID: "0", ActionID: "000"},

		events.ActionStrategyMarkedSequentialV1{OutcomeID: "9", ActionID: "99"},
		events.ActionStrategyMarkedSequentialV1{OutcomeID: "9", ActionID: "999"},

		events.OutcomeFixedV1{OutcomeID: "0"},
		events.OutcomeFixedV1{OutcomeID: "9"},

		events.OutcomeDeferredV1{OutcomeID: "1"},
		events.OutcomeDeferredV1{OutcomeID: "8"},

		events.OutcomeUncertainV1{OutcomeID: "2"},
		events.OutcomeUncertainV1{OutcomeID: "7"},

		events.OutcomeAbandonedV1{OutcomeID: "3"},
		events.OutcomeAbandonedV1{OutcomeID: "6"},

		events.OutcomeDeletedV1{OutcomeID: "4"},
		events.OutcomeRealizedV1{OutcomeID: "5"},
	)

	this.assert(IncompleteActionsByContext{
		Contexts: []*Context{
			{
				Name: "c0",
				Actions: []*ContextualAction{
					{
						ActionDetails: &ActionDetails{
							ID:          "00",
							Description: "action00",
							Contexts:    []string{"c0", "c1"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategyConcurrent,
						},
						OutcomeTitle: "0",
					},
					{
						ActionDetails: &ActionDetails{
							ID:          "99",
							Description: "action99",
							Contexts:    []string{"c0", "c1"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategySequential,
						},
						OutcomeTitle: "9",
					},
				},
			},
			{
				Name: "c1",
				Actions: []*ContextualAction{
					{
						ActionDetails: &ActionDetails{
							ID:          "00",
							Description: "action00",
							Contexts:    []string{"c0", "c1"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategyConcurrent,
						},
						OutcomeTitle: "0",
					},
					{
						ActionDetails: &ActionDetails{
							ID:          "99",
							Description: "action99",
							Contexts:    []string{"c0", "c1"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategySequential,
						},
						OutcomeTitle: "9",
					},
				},
			},
			{
				Name: "c2",
				Actions: []*ContextualAction{
					{
						ActionDetails: &ActionDetails{
							ID:          "0000",
							Description: "action0000",
							Contexts:    []string{"c2", "c3"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategyConcurrent,
						},
						OutcomeTitle: "0",
					},
				},
			},
			{
				Name: "c3",
				Actions: []*ContextualAction{
					{
						ActionDetails: &ActionDetails{
							ID:          "0000",
							Description: "action0000",
							Contexts:    []string{"c2", "c3"},
							Status:      core.ActionStatusIncomplete,
							Strategy:    core.ActionStrategyConcurrent,
						},
						OutcomeTitle: "0",
					},
				},
			},
		},
	})

}
