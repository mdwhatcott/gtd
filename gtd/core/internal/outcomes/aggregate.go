package outcomes

import (
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	log *logging.Logger
	now time.Time

	id          string
	title       string
	explanation string
	description string
	status      string
	deleted     bool

	results []interface{}
}

func NewAggregate(_now time.Time, _log *logging.Logger) *Aggregate {
	return &Aggregate{
		now: _now,
		log: _log,
	}
}

func (this *Aggregate) TrackOutcome(_outcomeID, _title string) error {
	return this.raise(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: _outcomeID,
			Title:     _title,
		},
		events.OutcomeFixedV1{
			Timestamp: this.now,
			OutcomeID: _outcomeID,
		},
	)
}
func (this *Aggregate) UpdateOutcomeTitle(_title string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if _title == this.title {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeTitleUpdatedV1{
		Timestamp:    this.now,
		OutcomeID:    this.id,
		UpdatedTitle: _title,
	})
}
func (this *Aggregate) UpdateOutcomeExplanation(_explanation string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if _explanation == this.explanation {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeExplanationUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		UpdatedExplanation: _explanation,
	})
}
func (this *Aggregate) UpdateOutcomeDescription(_description string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if _description == this.description {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeDescriptionUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		UpdatedDescription: _description,
	})
}
func (this *Aggregate) DeleteOutcome() error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if this.deleted {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeDeletedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeRealized() error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if this.status == "REALIZED" {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeRealizedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeFixed() error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if this.status == "FIXED" {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeFixedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeAbandoned() error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if this.status == "ABANDONED" {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeAbandonedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeDeferred() error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if this.status == "DEFERRED" {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeDeferredV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}

func (this *Aggregate) apply(_event interface{}) {
	switch EVENT := _event.(type) {

	case events.OutcomeTrackedV1:
		this.id = EVENT.OutcomeID
		this.title = EVENT.Title

	case events.OutcomeFixedV1:
		this.status = "FIXED"

	case events.OutcomeRealizedV1:
		this.status = "REALIZED"

	case events.OutcomeAbandonedV1:
		this.status = "ABANDONED"

	case events.OutcomeDeferredV1:
		this.status = "DEFERRED"

	case events.OutcomeTitleUpdatedV1:
		this.title = EVENT.UpdatedTitle

	case events.OutcomeExplanationUpdatedV1:
		this.explanation = EVENT.UpdatedExplanation

	case events.OutcomeDescriptionUpdatedV1:
		this.description = EVENT.UpdatedDescription

	case events.OutcomeDeletedV1:
		this.deleted = true
	}
}
func (this *Aggregate) raise(_events ...interface{}) error {
	for _, EVENT := range _events {
		this.results = append(this.results, EVENT)
		this.apply(EVENT)
	}
	return nil
}
func (this *Aggregate) Replay(events ...interface{}) {
	this.log.Println("stream:", len(events))
	for _, EVENT := range events {
		this.log.Println("applying event:", EVENT)
		this.apply(EVENT)
	}
}
func (this *Aggregate) TransferResults() []interface{} {
	RESULTS := this.results
	this.results = nil
	return RESULTS
}
