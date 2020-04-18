package domain

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	now time.Time

	id         string
	definition string

	results []interface{}
}

func NewAggregate(now time.Time) *Aggregate {
	return &Aggregate{now: now}
}
func (this *Aggregate) DefineOutcome(outcomeID, definition string) error {
	return this.raise(events.OutcomeExplanationProvidedV1{
		Timestamp:   this.now,
		OutcomeID:   outcomeID,
		Explanation: definition,
	})
}
func (this *Aggregate) RedefineOutcome(definition string) error {
	if len(this.id) == 0 {
		return core.ErrOutcomeNotFound
	}
	if definition == this.definition {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeRedefinedV1{
		Timestamp:     this.now,
		OutcomeID:     this.id,
		NewDefinition: definition,
	})
}
func (this *Aggregate) raise(event interface{}) error {
	this.results = append(this.results, event)
	return nil
}
func (this *Aggregate) apply(event interface{}) {
	switch event := event.(type) {
	case events.OutcomeExplanationProvidedV1:
		this.applyOutcomeDefined(event)
	}
}
func (this *Aggregate) applyOutcomeDefined(event events.OutcomeExplanationProvidedV1) {
	this.id = event.OutcomeID
	this.definition = event.Explanation
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
