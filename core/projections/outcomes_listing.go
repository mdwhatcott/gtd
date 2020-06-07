package projections

import (
	"sort"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/events"
)

type OutcomesListingProjector struct {
	OutcomesListing

	all map[string]*OutcomesListingItem
}

func NewOutcomesListingProjector() *OutcomesListingProjector {
	return &OutcomesListingProjector{
		all: make(map[string]*OutcomesListingItem),
	}
}

func (this *OutcomesListingProjector) Projection() interface{} {
	return this.OutcomesListingProjection()
}
func (this *OutcomesListingProjector) OutcomesListingProjection() OutcomesListing {
	return this.OutcomesListing
}

func (this *OutcomesListingProjector) Apply(_messages ...interface{}) {
	defer this.buildProjection()

	for _, MESSAGE := range _messages {
		switch EVENT := MESSAGE.(type) {
		case events.OutcomeTrackedV1:
			this.all[EVENT.OutcomeID] = &OutcomesListingItem{
				ID:     EVENT.OutcomeID,
				Title:  EVENT.Title,
				Status: core.OutcomeStatusFixed,
			}
		case events.OutcomeTitleUpdatedV1:
			this.all[EVENT.OutcomeID].Title = EVENT.UpdatedTitle
		case events.OutcomeFixedV1:
			this.all[EVENT.OutcomeID].Status = core.OutcomeStatusFixed
		case events.OutcomeDeferredV1:
			this.all[EVENT.OutcomeID].Status = core.OutcomeStatusDeferred
		case events.OutcomeUncertainV1:
			this.all[EVENT.OutcomeID].Status = core.OutcomeStatusUncertain
		case events.OutcomeAbandonedV1:
			this.all[EVENT.OutcomeID].Status = core.OutcomeStatusAbandoned
		case events.OutcomeRealizedV1:
			this.all[EVENT.OutcomeID].Status = core.OutcomeStatusRealized
		case events.OutcomeDeletedV1:
			delete(this.all, EVENT.OutcomeID)
		}
	}
}

func (this *OutcomesListingProjector) buildProjection() {
	this.Fixed = sortByTitle(this.filterByStatus(core.OutcomeStatusFixed))
	this.Deferred = sortByTitle(this.filterByStatus(core.OutcomeStatusDeferred))
	this.Uncertain = sortByTitle(this.filterByStatus(core.OutcomeStatusUncertain))
	this.Abandoned = sortByTitle(this.filterByStatus(core.OutcomeStatusAbandoned))
	this.Realized = sortByTitle(this.filterByStatus(core.OutcomeStatusRealized))
}

func (this *OutcomesListingProjector) filterByStatus(status core.OutcomeStatus) (filtered []*OutcomesListingItem) {
	for _, item := range this.all {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func sortByTitle(listing []*OutcomesListingItem) []*OutcomesListingItem {
	sort.Slice(listing, func(i, j int) bool {
		return listing[i].Title < listing[j].Title
	})
	return listing
}
