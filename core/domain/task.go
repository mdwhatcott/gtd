package domain

import (
	"context"
	"reflect"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/storage"
)

type Task struct {
	*joyride.Base

	log    core.Logger
	clock  core.Clock
	nextID core.IDFunc

	queries      map[string]*storage.OutcomeEventStream
	instructions []interface{}
	aggregates   map[string]*Aggregate
}

func NewTask(log core.Logger, clock core.Clock, nextID core.IDFunc) *Task {
	return &Task{
		Base:       joyride.New(),
		log:        log,
		clock:      clock,
		nextID:     nextID,
		queries:    make(map[string]*storage.OutcomeEventStream),
		aggregates: make(map[string]*Aggregate),
	}
}

func (this *Task) saveInstruction(command interface{}) {
	this.instructions = append(this.instructions, command)
}
func (this *Task) aggregate(id commands.Identifiable) *Aggregate {
	AGGREGATE, FOUND := this.aggregates[id.ID()]
	if FOUND {
		return AGGREGATE
	}
	return this.createAggregate(id.ID())
}
func (this *Task) createAggregate(id string) (aggregate_ *Aggregate) {
	aggregate_ = NewAggregate(this.clock(), this.log)
	this.aggregates[id] = aggregate_
	return aggregate_
}

func (this *Task) PrepareToTrackOutcome(command *commands.TrackOutcome) {
	this.saveInstruction(command)
}
func (this *Task) PrepareInstruction(instruction commands.Identifiable) {
	this.saveInstruction(instruction)
	this.registerOutcomeEventStreamQuery(instruction.ID())
}
func (this *Task) registerOutcomeEventStreamQuery(id string) {
	QUERY, FOUND := this.queries[id]
	if FOUND {
		return
	}
	QUERY = &storage.OutcomeEventStream{OutcomeID: id}
	this.queries[id] = QUERY
	this.AddRequiredReads(QUERY)
}

func (this *Task) Execute(_ context.Context) joyride.TaskResult {
	this.replayEvents()
	this.processInstructions()
	this.publishResults()
	return this
}
func (this *Task) replayEvents() {
	for _, QUERY := range this.queries {
		this.aggregate(QUERY).Replay(QUERY.Result.Events)
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

		default:
			this.log.Println("Unrecognized instruction:", reflect.TypeOf(COMMAND).String())
			continue
		}
	}
}
func (this *Task) publishResults() {
	for _, AGGREGATE := range this.aggregates {
		RESULTS := AGGREGATE.Results()
		for _, result := range RESULTS {
			this.AddPendingWrites(result)
		}
	}
}
