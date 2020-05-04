package outcomes

import (
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	now time.Time
	log *logging.Logger

	id          string
	title       string
	explanation string
	description string
	results     []interface{}
}

func NewAggregate(_now time.Time, _log *logging.Logger) *Aggregate {
	return &Aggregate{now: _now, log: _log}
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

func (this *Aggregate) apply(_event interface{}) {
	switch EVENT := _event.(type) {
	case events.OutcomeTrackedV1:
		this.id = EVENT.OutcomeID
		this.title = EVENT.Title
	case events.OutcomeTitleUpdatedV1:
		this.title = EVENT.UpdatedTitle
	case events.OutcomeExplanationUpdatedV1:
		this.explanation = EVENT.UpdatedExplanation
	case events.OutcomeDescriptionUpdatedV1:
		this.description = EVENT.UpdatedDescription
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
