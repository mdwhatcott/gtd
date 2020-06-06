package projections

import (
	"testing"

	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/events"
)

func TestOutcomesListingFixture(t *testing.T) {
	gunit.Run(new(OutcomesListingFixture), t)
}

type OutcomesListingFixture struct {
	*gunit.Fixture
	*ProjectorFixture
}

func (this *OutcomesListingFixture) Setup() {
	this.ProjectorFixture = InitializeProjectorFixture(this.Fixture, NewOutcomesListingProjector())
}

func (this *OutcomesListingFixture) TestOutcomeTracked() {
	this.apply(events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"})
	this.assert(OutcomesListing{
		Fixed: []*OutcomesListingItem{{ID: "1", Title: "title", Status: core.OutcomeStatusFixed}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeTitleUpdated() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeTitleUpdatedV1{OutcomeID: "1", UpdatedTitle: "updated-title"},
	)
	this.assert(OutcomesListing{
		Fixed: []*OutcomesListingItem{{ID: "1", Title: "updated-title", Status: core.OutcomeStatusFixed}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeFixed() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeDeferredV1{OutcomeID: "1"},
		events.OutcomeFixedV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{
		Fixed: []*OutcomesListingItem{{ID: "1", Title: "title", Status: core.OutcomeStatusFixed}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeDeferred() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeDeferredV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{
		Deferred: []*OutcomesListingItem{{ID: "1", Title: "title", Status: core.OutcomeStatusDeferred}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeUncertain() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeUncertainV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{
		Uncertain: []*OutcomesListingItem{{ID: "1", Title: "title", Status: core.OutcomeStatusUncertain}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeAbandoned() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeAbandonedV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{
		Abandoned: []*OutcomesListingItem{{ID: "1", Title: "title", Status: core.OutcomeStatusAbandoned}},
	})
}
func (this *OutcomesListingFixture) TestOutcomeRealized() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeRealizedV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{})
}
func (this *OutcomesListingFixture) TestOutcomeDeleted() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "title"},
		events.OutcomeDeletedV1{OutcomeID: "1"},
	)
	this.assert(OutcomesListing{})
}
func (this *OutcomesListingFixture) TestSortingOfListings() {
	this.apply(
		events.OutcomeTrackedV1{OutcomeID: "0", Title: "j"},
		events.OutcomeTrackedV1{OutcomeID: "1", Title: "i"},
		events.OutcomeTrackedV1{OutcomeID: "2", Title: "h"},
		events.OutcomeTrackedV1{OutcomeID: "3", Title: "g"},
		events.OutcomeTrackedV1{OutcomeID: "4", Title: "f"},
		events.OutcomeTrackedV1{OutcomeID: "5", Title: "e"},
		events.OutcomeTrackedV1{OutcomeID: "6", Title: "d"},
		events.OutcomeTrackedV1{OutcomeID: "7", Title: "c"},
		events.OutcomeTrackedV1{OutcomeID: "8", Title: "b"},
		events.OutcomeTrackedV1{OutcomeID: "9", Title: "a"},

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
	this.assert(OutcomesListing{
		Fixed: []*OutcomesListingItem{
			{ID: "9", Title: "a", Status: core.OutcomeStatusFixed},
			{ID: "0", Title: "j", Status: core.OutcomeStatusFixed},
		},
		Deferred: []*OutcomesListingItem{
			{ID: "8", Title: "b", Status: core.OutcomeStatusDeferred},
			{ID: "1", Title: "i", Status: core.OutcomeStatusDeferred},
		},
		Uncertain: []*OutcomesListingItem{
			{ID: "7", Title: "c", Status: core.OutcomeStatusUncertain},
			{ID: "2", Title: "h", Status: core.OutcomeStatusUncertain},
		},
		Abandoned: []*OutcomesListingItem{
			{ID: "6", Title: "d", Status: core.OutcomeStatusAbandoned},
			{ID: "3", Title: "g", Status: core.OutcomeStatusAbandoned},
		},
	})
}
