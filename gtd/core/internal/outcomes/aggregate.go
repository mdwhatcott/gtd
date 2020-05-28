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
	status      core.OutcomeStatus
	deleted     bool
	actions     map[string]*Action

	results []interface{}
}

func NewAggregate(_now time.Time, _log *logging.Logger) *Aggregate {
	return &Aggregate{
		now:     _now,
		log:     _log,
		actions: make(map[string]*Action),
	}
}

func (this *Aggregate) TransferResults() []interface{} {
	RESULTS := this.results
	this.results = nil
	return RESULTS
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

func (this *Aggregate) exists() bool {
	if len(this.id) == 0 {
		return false
	}
	if this.deleted {
		return false
	}
	return true
}

func (this *Aggregate) apply(_event interface{}) {
	switch EVENT := _event.(type) {

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

func (this *Aggregate) TrackOutcome(_outcomeID, _title string) {
	_ = this.raise(
		events.OutcomeTrackedV1{
			Timestamp: this.now,
			OutcomeID: _outcomeID,
			Title:     _title,
		},
	)
}
func (this *Aggregate) UpdateOutcomeTitle(_title string) error {
	if !this.exists() {
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
	if !this.exists() {
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
	if !this.exists() {
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
func (this *Aggregate) TrackAction(_id, _description string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	return this.raise(events.ActionTrackedV1{
		Timestamp:   this.now,
		OutcomeID:   this.id,
		ActionID:    _id,
		Description: _description,
		Contexts:    gatherContexts(_description),
	})
}
func (this *Aggregate) UpdateActionDescription(_id, _description string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Description == _description {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionDescriptionUpdatedV1{
		Timestamp:          this.now,
		OutcomeID:          this.id,
		ActionID:           _id,
		UpdatedDescription: _description,
		UpdatedContexts:    gatherContexts(_description),
	})
}
func (this *Aggregate) ReorderActions(_newIDOrder []string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	if len(this.actions) == 0 {
		return core.ErrActionNotFound
	}
	if len(_newIDOrder) != len(this.actions) {
		return core.ErrActionNotFound
	}
	for _, ID := range _newIDOrder {
		if this.actions[ID] == nil {
			return core.ErrActionNotFound
		}
	}
	return this.raise(events.ActionsReorderedV1{
		Timestamp:  this.now,
		OutcomeID:  this.id,
		NewIDOrder: _newIDOrder,
	})
}
func (this *Aggregate) MarkActionStatusLatent(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Status == core.ActionStatusLatent {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedLatentV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
func (this *Aggregate) MarkActionStatusIncomplete(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Status == core.ActionStatusIncomplete {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedIncompleteV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
func (this *Aggregate) MarkActionStatusComplete(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Status == core.ActionStatusComplete {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStatusMarkedCompleteV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
func (this *Aggregate) MarkActionStrategySequential(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Strategy == core.ActionStrategySequential {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStrategyMarkedSequentialV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
func (this *Aggregate) MarkActionStrategyConcurrent(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	if action.Strategy == core.ActionStrategyConcurrent {
		return core.ErrOutcomeUnchanged
	}
	return this.raise(events.ActionStrategyMarkedConcurrentV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
func (this *Aggregate) DeleteAction(_id string) error {
	if !this.exists() {
		return core.ErrOutcomeNotFound
	}
	action := this.actions[_id]
	if action == nil {
		return core.ErrActionNotFound
	}
	return this.raise(events.ActionDeletedV1{
		Timestamp: this.now,
		OutcomeID: this.id,
		ActionID:  _id,
	})
}
