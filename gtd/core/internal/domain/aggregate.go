package domain

import (
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	now time.Time
	log *logging.Logger

	id    string
	title string

	results []interface{}
}

func NewAggregate(now time.Time, log *logging.Logger) *Aggregate {
	return &Aggregate{now: now, log: log}
}
func (this *Aggregate) TrackOutcome(outcomeID, title string) error {
	return this.raise(events.OutcomeTrackedV1{
		Timestamp: this.now,
		OutcomeID: outcomeID,
		Title:     title,
	})
}
func (this *Aggregate) ProvideOutcomeExplanation(explanation string) error {
	return this.raise(events.OutcomeExplanationProvidedV1{
		Timestamp:   this.now,
		OutcomeID:   this.id,
		Explanation: explanation,
	})
}
func (this *Aggregate) UpdateOutcomeTitle(title string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	return this.raise(events.OutcomeTitleUpdatedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		NewTitle:  title,
	})
}
func (this *Aggregate) raise(event interface{}) error {
	this.results = append(this.results, event)
	this.apply(event)
	return nil
}
func (this *Aggregate) apply(event interface{}) {
	switch event := event.(type) {
	case events.OutcomeTrackedV1:
		this.id = event.OutcomeID
	}
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
