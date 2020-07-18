package domain

import (
	"time"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/events"
	"github.com/mdwhatcott/gtd/v3/util/errors"
)

type Aggregate struct {
	log *logging.Logger
	now time.Time

	id          string
	title       string
	explanation string
	description string
	status      core.OutcomeStatus
	deleted     bool
	actions     map[string]*Action

	results []interface{}
}

func NewAggregate(now time.Time, log *logging.Logger) *Aggregate {
	return &Aggregate{
		now:     now,
		log:     log,
		actions: make(map[string]*Action),
	}
}

func (this *Aggregate) Results() []interface{} {
	return this.results
}
func (this *Aggregate) raise(events ...interface{}) error {
	for _, EVENT := range events {
		this.results = append(this.results, EVENT)
		this.apply(EVENT)
	}
	return nil
}
func (this *Aggregate) Replay(events chan interface{}) {
	for EVENT := range events {
		this.apply(EVENT)
	}
}

func (this *Aggregate) exists() bool {
	if len(this.id) == 0 {
		return false
	}
	if this.deleted {
		return false
	}
	return true
}

func (this *Aggregate) apply(event interface{}) {
	switch EVENT := event.(type) {

	case events.OutcomeTrackedV1:
		this.id = EVENT.OutcomeID
		this.title = EVENT.Title

	case events.OutcomeFixedV1:
		this.status = core.OutcomeStatusFixed

	case events.OutcomeRealizedV1:
		this.status = core.OutcomeStatusRealized

	case events.OutcomeAbandonedV1:
		this.status = core.OutcomeStatusAbandoned

	case events.OutcomeDeferredV1:
		this.status = core.OutcomeStatusDeferred

	case events.OutcomeUncertainV1:
		this.status = core.OutcomeStatusUncertain

	case events.OutcomeTitleUpdatedV1:
		this.title = EVENT.UpdatedTitle

	case events.OutcomeExplanationUpdatedV1:
		this.explanation = EVENT.UpdatedExplanation

	case events.OutcomeDescriptionUpdatedV1:
		this.description = EVENT.UpdatedDescription

	case events.OutcomeDeletedV1:
		this.deleted = true

	case events.ActionTrackedV1:
		this.actions[EVENT.ActionID] = &Action{Description: EVENT.Description}

	case events.ActionDescriptionUpdatedV1:
		this.actions[EVENT.ActionID].Description = EVENT.UpdatedDescription

	case events.ActionStatusMarkedLatentV1:
		this.actions[EVENT.ActionID].Status = core.ActionStatusLatent

	case events.ActionStatusMarkedIncompleteV1:
		this.actions[EVENT.ActionID].Status = core.ActionStatusIncomplete

	case events.ActionStatusMarkedCompleteV1:
		this.actions[EVENT.ActionID].Status = core.ActionStatusComplete

	case events.ActionStrategyMarkedSequentialV1:
		this.actions[EVENT.ActionID].Strategy = core.ActionStrategySequential

	case events.ActionStrategyMarkedConcurrentV1:
		this.actions[EVENT.ActionID].Strategy = core.ActionStrategyConcurrent

	case events.ActionDeletedV1:
		delete(this.actions, EVENT.ActionID)
	}
}

func (this *Aggregate) TrackOutcome(outcomeID, title string) {
	_ = this.raise(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: outcomeID,
			Title:     title,
		},
	)
}
func (this *Aggregate) UpdateOutcomeTitle(title string) error {
	if !this.exists() {
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
	if !this.exists() {
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
	if !this.exists() {
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
func (this *Aggregate) DeleteOutcome() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	return this.raise(events.OutcomeDeletedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeRealized() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if this.status == core.OutcomeStatusRealized {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeRealizedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeFixed() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if this.status == core.OutcomeStatusFixed {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeFixedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeAbandoned() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if this.status == core.OutcomeStatusAbandoned {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeAbandonedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeDeferred() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if this.status == core.OutcomeStatusDeferred {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeDeferredV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) DeclareOutcomeUncertain() error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if this.status == core.OutcomeStatusUncertain {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.OutcomeUncertainV1{
		Timestamp: this.now,
		OutcomeID: this.id,
	})
}
func (this *Aggregate) TrackAction(id, description string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	return this.raise(events.ActionTrackedV1{
		Timestamp:   this.now,
		OutcomeID:   this.id,
		ActionID:    id,
		Description: description,
		Contexts:    gatherContexts(description),
	})
}
func (this *Aggregate) UpdateActionDescription(id, description string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Description == description {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionDescriptionUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		ActionID:           id,
		UpdatedDescription: description,
		UpdatedContexts:    gatherContexts(description),
	})
}
func (this *Aggregate) ReorderActions(newIDOrder []string) error {
	// TODO: it should be OK if the newIDOrder has extra entries (but not the other way around)

	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if len(this.actions) == 0 {
		return errors.Wrap(core.ErrActionNotFound, errors.New("there are none"))
	}
	if len(newIDOrder) != len(this.actions) {
		return core.ErrActionNotFound
	}
	for _, ID := range newIDOrder {
		if this.actions[ID] == nil {
			return core.ErrActionNotFound
		}
	}
	return this.raise(events.ActionsReorderedV1{
		Timestamp:    this.now,
		OutcomeID:    this.id,
		ReorderedIDs: newIDOrder,
	})
}
func (this *Aggregate) MarkActionStatusLatent(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Status == core.ActionStatusLatent {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedLatentV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
func (this *Aggregate) MarkActionStatusIncomplete(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Status == core.ActionStatusIncomplete {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedIncompleteV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
func (this *Aggregate) MarkActionStatusComplete(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Status == core.ActionStatusComplete {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedCompleteV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
func (this *Aggregate) MarkActionStrategySequential(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Strategy == core.ActionStrategySequential {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStrategyMarkedSequentialV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
func (this *Aggregate) MarkActionStrategyConcurrent(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	if ACTION.Strategy == core.ActionStrategyConcurrent {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStrategyMarkedConcurrentV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
func (this *Aggregate) DeleteAction(id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	ACTION := this.actions[id]
	if ACTION == nil {
		return core.ErrActionNotFound
	}
	return this.raise(events.ActionDeletedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  id,
	})
}
