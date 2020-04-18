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
func (this *Task) DefineOutcome(message *commands.DefineOutcome) {
	this.instructions = append(this.instructions, message)
}
func (this *Task) RedefineOutcome(message *commands.RedefineOutcome) {
	this.instructions = append(this.instructions, message)
	this.registerOutcomeEventStreamQuery(message.OutcomeID)
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
		switch message := message.(type) {
		case *commands.DefineOutcome:
			this.trackOutcome(message)
		case *commands.RedefineOutcome:
			this.redefineOutcome(message)
		}
	}
}
func (this *Task) trackOutcome(command *commands.DefineOutcome) {
	outcomeID := this.nextID()
	aggregate := this.aggregate(outcomeID)
	command.Result.OutcomeID = outcomeID
	command.Result.Error = aggregate.DefineOutcome(outcomeID, command.Definition)
}
func (this *Task) redefineOutcome(command *commands.RedefineOutcome) {
	aggregate := this.aggregate(command.OutcomeID)
	command.Result.Error = aggregate.RedefineOutcome(command.NewDefinition)
}
func (this *Task) publishResults() {
	for _, aggregate := range this.aggregates {
		this.AddPendingWrites(aggregate.TransferResults()...)
	}
}
