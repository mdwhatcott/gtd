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
		aggregate = NewAggregate(this.clock.UTCNow())
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
	query.Result.Stream = make(chan interface{})
	this.queries[id] = query
	this.AddRequiredReads(query)
}
func (this *Task) TrackOutcome(command *commands.TrackOutcome) {
	id := this.nextID()
	aggregate := this.aggregate(id)
	command.Result.ID = id
	command.Result.Error = aggregate.TrackOutcome(id, command.Title)
}
func (this *Task) Execute() joyride.TaskResult {
	this.replayEvents()
	this.processInstructions()
	this.publishResults()
	return this
}
func (this *Task) replayEvents() {
	for id, query := range this.queries {
		this.aggregate(id).Replay(query.Result.Stream)
	}
}
func (this *Task) processInstructions() {
	for _, message := range this.instructions {
		switch message.(type) {
		}
	}
}
func (this *Task) trackOutcome(command *commands.TrackOutcome) {
	outcomeID := this.nextID()
	aggregate := this.aggregate(outcomeID)
	command.Result.ID = outcomeID
	command.Result.Error = aggregate.TrackOutcome(outcomeID, command.Title)
}
func (this *Task) publishResults() {
	for _, aggregate := range this.aggregates {
		this.AddPendingWrites(aggregate.TransferResults()...)
	}
}
