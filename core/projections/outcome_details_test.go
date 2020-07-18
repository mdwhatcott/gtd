package projections

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core"

	"github.com/mdwhatcott/gtd/v3/core/events"
)

func TestOutcomeDetailsFixture(t *testing.T) {
	gunit.Run(new(OutcomeDetailsFixture), t)
}

type OutcomeDetailsFixture struct {
	*gunit.Fixture
	*ProjectorFixture
}

func (this *OutcomeDetailsFixture) Setup() {
	this.ProjectorFixture = InitializeProjectorFixture(this.Fixture, NewOutcomeDetailsProjector())
}
func (this *OutcomeDetailsFixture) TestBlankWhenFirstInstantiated() {
	this.assert(OutcomeDetails{})
}
func (this *OutcomeDetailsFixture) TestOutcomeTracked() {
	this.apply(events.OutcomeTrackedV1{OutcomeID: "ID", Title: "title"})
	this.assert(OutcomeDetails{ID: "ID", Title: "title"})
}
func (this *OutcomeDetailsFixture) TestOutcomeTitleUpdated() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.OutcomeTitleUpdatedV1{UpdatedTitle: "updated-title"},
	)
	this.assert(OutcomeDetails{Title: "updated-title"})
}
func (this *OutcomeDetailsFixture) TestOutcomeDescriptionUpdated() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.OutcomeDescriptionUpdatedV1{UpdatedDescription: "description"},
	)
	this.assert(OutcomeDetails{Title: "title", Description: "description"})
}
func (this *OutcomeDetailsFixture) TestOutcomeExplanationUpdated() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.OutcomeExplanationUpdatedV1{UpdatedExplanation: "explanation"},
	)
	this.assert(OutcomeDetails{Title: "title", Explanation: "explanation"})
}
func (this *OutcomeDetailsFixture) TestOutcomeFixed() {
	this.apply(
		events.OutcomeTrackedV1{},
		events.OutcomeFixedV1{},
	)
	this.assert(OutcomeDetails{Status: core.OutcomeStatusFixed})
}
func (this *OutcomeDetailsFixture) TestOutcomeDeferred() {
	this.apply(
		events.OutcomeTrackedV1{},
		events.OutcomeDeferredV1{},
	)
	this.assert(OutcomeDetails{Status: core.OutcomeStatusDeferred})
}
func (this *OutcomeDetailsFixture) TestOutcomeUncertain() {
	this.apply(
		events.OutcomeTrackedV1{},
		events.OutcomeUncertainV1{},
	)
	this.assert(OutcomeDetails{Status: core.OutcomeStatusUncertain})
}
func (this *OutcomeDetailsFixture) TestOutcomeAbandoned() {
	this.apply(
		events.OutcomeTrackedV1{},
		events.OutcomeAbandonedV1{},
	)
	this.assert(OutcomeDetails{Status: core.OutcomeStatusAbandoned})
}
func (this *OutcomeDetailsFixture) TestActionTracked() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionDeleted() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionDeletedV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionDeleted_NotPreviouslyTracked_Panic() {
	action := func() { this.apply(events.ActionDeletedV1{ActionID: "not-found"}) }
	this.So(action, should.Panic)
}
func (this *OutcomeDetailsFixture) TestActionDescriptionUpdated() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionDescriptionUpdatedV1{
			ActionID:           "1",
			UpdatedDescription: "updated-description",
			UpdatedContexts:    []string{"updated", "contexts"},
		},
	)
	this.assert(OutcomeDetails{
		Title: "title",
		Actions: []*ActionDetails{
			{
				ID:          "0",
				Description: "action-description0",
				Contexts:    []string{"context1", "context2"},
				Status:      core.ActionStatusIncomplete,
				Strategy:    core.ActionStrategyConcurrent,
			},
			{
				ID:          "1",
				Description: "updated-description",
				Contexts:    []string{"updated", "contexts"},
				Status:      core.ActionStatusIncomplete,
				Strategy:    core.ActionStrategyConcurrent,
			},
		},
	})
}
func (this *OutcomeDetailsFixture) TestActionStatusMarkedIncomplete() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionStatusMarkedCompleteV1{ActionID: "1"},
		events.ActionStatusMarkedIncompleteV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionStatusMarkedComplete() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionStatusMarkedCompleteV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusComplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionStatusMarkedLatent() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionStatusMarkedLatentV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusLatent,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionStrategyMarkedSequential() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionStrategyMarkedSequentialV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategySequential,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionStrategyMarkedConcurrent() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
		},
		events.ActionStrategyMarkedSequentialV1{ActionID: "1"},
		events.ActionStrategyMarkedConcurrentV1{ActionID: "1"},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "0",
					Description: "action-description0",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
func (this *OutcomeDetailsFixture) TestActionsReordered() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
		},
		events.ActionTrackedV1{
			ActionID:    "2",
			Description: "action-description2",
		},
		events.ActionsReorderedV1{
			ReorderedIDs: []string{"2", "0", "1"},
		},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "2",
					Description: "action-description2",
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "0",
					Description: "action-description0",
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Status:      core.ActionStatusIncomplete,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
