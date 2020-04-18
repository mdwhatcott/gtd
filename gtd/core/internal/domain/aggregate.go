package domain

import (
	"errors"
	"time"

	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Aggregate struct {
	now time.Time

	id         string
	definition string

	results []interface{}
}

func NewAggregate(id string, now time.Time) *Aggregate {
	return &Aggregate{id: id, now: now}
}
func (this *Aggregate) DefineOutcome(definition string) {
	_ = this.raise(events.OutcomeDefinedV1{
		Timestamp:  this.now,
		OutcomeID:  this.id,
		Definition: definition,
	})
}
func (this *Aggregate) RedefineOutcome(definition string) error {
	if definition == this.definition {
		return errors.New("HI") // TODO: more specific, application/contractual value
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
	case events.OutcomeDefinedV1:
		this.applyOutcomeDefined(event)
	}
}
func (this *Aggregate) applyOutcomeDefined(event events.OutcomeDefinedV1) {
	this.id = event.OutcomeID
	this.definition = event.Definition
}
func (this *Aggregate) Replay(stream chan interface{}) {
	for event := range stream {
		this.apply(event)
	}
}
func (this *Aggregate) TransferResults() []interface{} {
	return this.results // TODO: clear out this.results
}
