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

func NewAggregate(now time.Time, log *logging.Logger) *Aggregate {
	return &Aggregate{now: now, log: log}
}

func (this *Aggregate) TrackOutcome(outcomeID, title string) error {
	return this.raise(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: outcomeID,
			Title:     title,
		},
		events.OutcomeFixedV1{
			Timestamp: this.now,
			OutcomeID: outcomeID,
		},
	)
}

func (this *Aggregate) UpdateOutcomeTitle(title string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if title == this.title {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeTitleUpdatedV1{
		Timestamp:    this.now,
		OutcomeID:    this.id,
		UpdatedTitle: title,
	})
}

func (this *Aggregate) UpdateOutcomeExplanation(explanation string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if explanation == this.explanation {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeExplanationUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		UpdatedExplanation: explanation,
	})
}

func (this *Aggregate) UpdateOutcomeDescription(description string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if description == this.description {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeDescriptionUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		UpdatedDescription: description,
	})
}

func (this *Aggregate) apply(event interface{}) {
	switch event := event.(type) {
	case events.OutcomeTrackedV1:
		this.id = event.OutcomeID
		this.title = event.Title
	case events.OutcomeTitleUpdatedV1:
		this.title = event.UpdatedTitle
	case events.OutcomeExplanationUpdatedV1:
		this.explanation = event.UpdatedExplanation
	case events.OutcomeDescriptionUpdatedV1:
		this.description = event.UpdatedDescription
	}
}

func (this *Aggregate) raise(events ...interface{}) error {
	for _, event := range events {
		this.results = append(this.results, event)
		this.apply(event)
	}
	return nil
}

func (this *Aggregate) Replay(stream chan interface{}) {
	this.log.Println("stream:", len(stream))
	for event := range stream {
		this.log.Println("applying event:", event)
		this.apply(event)
	}
}

func (this *Aggregate) TransferResults() []interface{} {
	results := this.results
	this.results = nil
	return results
}
