package projections

import (
	"sort"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
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
			delete(this.all, EVENT.OutcomeID)
		case events.OutcomeDeletedV1:
			delete(this.all, EVENT.OutcomeID)
		}
	}
}

func (this *OutcomesListingProjector) buildProjection() {
	this.Fixed = this.filter(core.OutcomeStatusFixed)
	this.Deferred = this.filter(core.OutcomeStatusDeferred)
	this.Uncertain = this.filter(core.OutcomeStatusUncertain)
	this.Abandoned = this.filter(core.OutcomeStatusAbandoned)
}

func (this *OutcomesListingProjector) filter(status core.OutcomeStatus) (filtered []*OutcomesListingItem) {
	for _, item := range this.all {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})
	return filtered
}
