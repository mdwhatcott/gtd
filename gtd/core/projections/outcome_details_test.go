package projections

import (
	"testing"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
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
	this.apply(events.OutcomeTrackedV1{Title: "title"})
	this.assert(OutcomeDetails{Title: "title"})
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
func (this *OutcomeDetailsFixture) TestActionTracked() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    0,
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
					Status:      core.ActionStatusLatent,
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
func (this *OutcomeDetailsFixture) TestActionDeleted() {
	this.apply(
		events.OutcomeTrackedV1{Title: "title"},
		events.ActionTrackedV1{
			ActionID:    "0",
			Description: "action-description0",
			Contexts:    []string{"context1", "context2"},
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    0,
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
					Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    0,
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
				Status:      core.ActionStatusLatent,
				Strategy:    core.ActionStrategyConcurrent,
			},
			{
				ID:          "1",
				Description: "updated-description",
				Contexts:    []string{"updated", "contexts"},
				Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    1,
		},
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
					Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    1,
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
					Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    1,
		},
		events.ActionStatusMarkedCompleteV1{ActionID: "1"},
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
					Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    1,
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
					Status:      core.ActionStatusLatent,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Contexts:    []string{"context1", "context2"},
					Status:      core.ActionStatusLatent,
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
			Sequence:    0,
		},
		events.ActionTrackedV1{
			ActionID:    "1",
			Description: "action-description1",
			Contexts:    []string{"context1", "context2"},
			Sequence:    1,
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
					Status:      core.ActionStatusLatent,
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
			NewIDOrder: []string{"2", "0", "1"},
		},
	)
	this.assert(
		OutcomeDetails{
			Title: "title",
			Actions: []*ActionDetails{
				{
					ID:          "2",
					Description: "action-description2",
					Status:      core.ActionStatusLatent,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "0",
					Description: "action-description0",
					Status:      core.ActionStatusLatent,
					Strategy:    core.ActionStrategyConcurrent,
				},
				{
					ID:          "1",
					Description: "action-description1",
					Status:      core.ActionStatusLatent,
					Strategy:    core.ActionStrategyConcurrent,
				},
			},
		},
	)
}
