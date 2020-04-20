package domain

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Task struct {
	*joyride.Base

	log    *logging.Logger
	clock  *clock.Clock
	nextID func() string

	queries      map[string]*storage.OutcomeEventStream
	instructions []interface{}
	aggregates   map[string]*Aggregate
}

func NewTask(nextID func() string) *Task {
	return &Task{
		Base:       joyride.New(),
		nextID:     nextID,
		queries:    make(map[string]*storage.OutcomeEventStream),
		aggregates: make(map[string]*Aggregate),
	}
}
func (this *Task) aggregate(id string) *Aggregate {
	aggregate, found := this.aggregates[id]
	if !found {
		aggregate = NewAggregate(this.clock.UTCNow(), this.log)
		this.aggregates[id] = aggregate
	}
	return aggregate
}
func (this *Task) registerOutcomeEventStreamQuery(id string) {
	query, found := this.queries[id]
	if found {
		return
	}
	query = &storage.OutcomeEventStream{OutcomeID: id}
	this.queries[id] = query
	this.AddRequiredReads(query)
}
func (this *Task) TrackOutcome(command *commands.TrackOutcome) {
	this.instructions = append(this.instructions, command)
}
func (this *Task) UpdateOutcomeTitle(command *commands.UpdateOutcomeTitle) {
	this.instructions = append(this.instructions, command)
	this.registerOutcomeEventStreamQuery(command.OutcomeID)
}
func (this *Task) UpdateOutcomeExplanation(command *commands.UpdateOutcomeExplanation) {
	this.instructions = append(this.instructions, command)
	this.registerOutcomeEventStreamQuery(command.OutcomeID)
}
func (this *Task) Execute() joyride.TaskResult {
	this.replayEvents()
	this.processInstructions()
	this.publishResults()
	return this
}
func (this *Task) replayEvents() {
	for id, query := range this.queries {
		this.aggregate(id).Replay(query.Result.Events)
	}
}
func (this *Task) processInstructions() {
	for _, message := range this.instructions {
		switch command := message.(type) {
		case *commands.TrackOutcome:
			this.trackOutcome(command)
		case *commands.UpdateOutcomeTitle:
			this.updateOutcomeTitle(command)
		case *commands.UpdateOutcomeExplanation:
			this.updateOutcomeExplanation(command)
		}
	}
}
func (this *Task) trackOutcome(command *commands.TrackOutcome) {
	outcomeID := this.nextID()
	aggregate := this.aggregate(outcomeID)
	command.Result.ID = outcomeID
	command.Result.Error = aggregate.TrackOutcome(outcomeID, command.Title)
}
func (this *Task) updateOutcomeTitle(command *commands.UpdateOutcomeTitle) {
	aggregate := this.aggregate(command.OutcomeID)
	command.Result.Error = aggregate.UpdateOutcomeTitle(command.NewTitle)
}
func (this *Task) updateOutcomeExplanation(command *commands.UpdateOutcomeExplanation) {
	aggregate := this.aggregate(command.OutcomeID)
	command.Result.Error = aggregate.UpdateOutcomeExplanation(command.Explanation)
}
func (this *Task) publishResults() {
	for _, aggregate := range this.aggregates {
		this.AddPendingWrites(aggregate.TransferResults()...)
	}
}
