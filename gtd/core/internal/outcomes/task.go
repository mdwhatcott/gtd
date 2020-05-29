package outcomes

import (
	"github.com/smartystreets/clock"
	"github.com/smartystreets/joyride/v2"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Task struct {
	*joyride.Base

	log    *logging.Logger
	clock  *clock.Clock
	nextID core.IDFunc

	queries      map[string]*storage.OutcomeEventStream
	instructions []interface{}
	aggregates   map[string]*Aggregate
}

func NewTask(_nextID core.IDFunc) *Task {
	return &Task{
		Base:       joyride.New(),
		nextID:     _nextID,
		queries:    make(map[string]*storage.OutcomeEventStream),
		aggregates: make(map[string]*Aggregate),
	}
}

func (this *Task) saveInstruction(_command interface{}) {
	this.instructions = append(this.instructions, _command)
}
func (this *Task) aggregate(_id commands.Identifiable) *Aggregate {
	AGGREGATE, FOUND := this.aggregates[_id.ID()]
	if FOUND {
		return AGGREGATE
	}
	return this.createAggregate(_id.ID())
}
func (this *Task) createAggregate(_id string) (aggregate_ *Aggregate) {
	aggregate_ = NewAggregate(this.clock.UTCNow(), this.log)
	this.aggregates[_id] = aggregate_
	return aggregate_
}

func (this *Task) PrepareToTrackOutcome(_command *commands.TrackOutcome) {
	this.saveInstruction(_command)
}
func (this *Task) PrepareInstruction(_instruction commands.Identifiable) {
	this.saveInstruction(_instruction)
	this.registerOutcomeEventStreamQuery(_instruction.ID())
}
func (this *Task) registerOutcomeEventStreamQuery(_id string) {
	QUERY, FOUND := this.queries[_id]
	if FOUND {
		return
	}
	QUERY = &storage.OutcomeEventStream{OutcomeID: _id}
	this.queries[_id] = QUERY
	this.AddRequiredReads(QUERY)
}

func (this *Task) Execute() joyride.TaskResult {
	this.replayEvents()
	this.processInstructions()
	this.publishResults()
	return this
}
func (this *Task) replayEvents() {
	for _, QUERY := range this.queries {
		this.aggregate(QUERY).Replay(QUERY.Result.Events...)
	}
}
func (this *Task) processInstructions() {
	for _, MESSAGE := range this.instructions {
		switch COMMAND := MESSAGE.(type) {

		case *commands.TrackOutcome:
			COMMAND.Result.ID = this.nextID()
			this.createAggregate(COMMAND.Result.ID).TrackOutcome(COMMAND.Result.ID, COMMAND.Title)

		case *commands.UpdateOutcomeTitle:
			COMMAND.Result.Error = this.aggregate(COMMAND).UpdateOutcomeTitle(COMMAND.UpdatedTitle)

		case *commands.UpdateOutcomeExplanation:
			COMMAND.Result.Error = this.aggregate(COMMAND).UpdateOutcomeExplanation(COMMAND.UpdatedExplanation)

		case *commands.UpdateOutcomeDescription:
			COMMAND.Result.Error = this.aggregate(COMMAND).UpdateOutcomeDescription(COMMAND.UpdatedDescription)

		case *commands.DeleteOutcome:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeleteOutcome()

		case *commands.DeclareOutcomeRealized:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeclareOutcomeRealized()

		case *commands.DeclareOutcomeFixed:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeclareOutcomeFixed()

		case *commands.DeclareOutcomeAbandoned:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeclareOutcomeAbandoned()

		case *commands.DeclareOutcomeDeferred:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeclareOutcomeDeferred()

		case *commands.DeclareOutcomeUncertain:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeclareOutcomeUncertain()

		case *commands.TrackAction:
			ID := this.nextID()
			COMMAND.Result.Error = this.aggregate(COMMAND).TrackAction(ID, COMMAND.Description)
			if COMMAND.Result.Error == nil {
				COMMAND.Result.ID = ID
			}

		case *commands.UpdateActionDescription:
			COMMAND.Result.Error = this.aggregate(COMMAND).UpdateActionDescription(COMMAND.ActionID, COMMAND.UpdatedDescription)

		case *commands.ReorderActions:
			COMMAND.Result.Error = this.aggregate(COMMAND).ReorderActions(COMMAND.ReorderedIDs)

		case *commands.MarkActionStatusLatent:
			COMMAND.Result.Error = this.aggregate(COMMAND).MarkActionStatusLatent(COMMAND.ActionID)

		case *commands.MarkActionStatusIncomplete:
			COMMAND.Result.Error = this.aggregate(COMMAND).MarkActionStatusIncomplete(COMMAND.ActionID)

		case *commands.MarkActionStatusComplete:
			COMMAND.Result.Error = this.aggregate(COMMAND).MarkActionStatusComplete(COMMAND.ActionID)

		case *commands.MarkActionStrategySequential:
			COMMAND.Result.Error = this.aggregate(COMMAND).MarkActionStrategySequential(COMMAND.ActionID)

		case *commands.MarkActionStrategyConcurrent:
			COMMAND.Result.Error = this.aggregate(COMMAND).MarkActionStrategyConcurrent(COMMAND.ActionID)

		case *commands.DeleteAction:
			COMMAND.Result.Error = this.aggregate(COMMAND).DeleteAction(COMMAND.ActionID)
		}
	}
}
func (this *Task) publishResults() {
	for _, AGGREGATE := range this.aggregates {
		this.AddPendingWrites(AGGREGATE.TransferResults()...)
	}
}
