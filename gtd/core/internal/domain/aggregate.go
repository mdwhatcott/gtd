package domain

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	now time.Time

	id    string
	title string

	results []interface{}
}

func NewAggregate(now time.Time) *Aggregate {
	return &Aggregate{now: now}
}
func (this *Aggregate) TrackOutcome(outcomeID, title string) error {
	return this.raise(events.OutcomeExplanationProvidedV1{
		Timestamp:   this.now,
		OutcomeID:   outcomeID,
		Explanation: title,
	})
}
func (this *Aggregate) raise(event interface{}) error {
	this.results = append(this.results, event)
	return nil
}
func (this *Aggregate) apply(event interface{}) {
	switch event := event.(type) {
	case events.OutcomeTrackedV1:
		this.applyOutcomeDefined(event)
	}
}
func (this *Aggregate) applyOutcomeDefined(event events.OutcomeTrackedV1) {
	this.id = event.OutcomeID
	this.title = event.Title
}
func (this *Aggregate) Replay(stream chan interface{}) {
	for event := range stream {
		this.apply(event)
	}
}
func (this *Aggregate) TransferResults() []interface{} {
	results := this.results
	this.results = nil
	return results
}
