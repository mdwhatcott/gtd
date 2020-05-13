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

func (this *Task) aggregate(_id string) *Aggregate {
	AGGREGATE, FOUND := this.aggregates[_id]
	if !FOUND {
		AGGREGATE = NewAggregate(this.clock.UTCNow(), this.log)
		this.aggregates[_id] = AGGREGATE
	}
	return AGGREGATE
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

func (this *Task) PrepareToTrackOutcome(_command *commands.TrackOutcome) {
	this.instructions = append(this.instructions, _command)
}

func (this *Task) PrepareInstruction(_instruction interface{}, _id string) {
	this.instructions = append(this.instructions, _instruction)
	this.registerOutcomeEventStreamQuery(_id)
}

func (this *Task) Execute() joyride.TaskResult {
	this.replayEvents()
	this.processInstructions()
	this.publishResults()
	return this
}

func (this *Task) replayEvents() {
	for ID, QUERY := range this.queries {
		this.aggregate(ID).Replay(QUERY.Result.Events...)
	}
}

func (this *Task) processInstructions() {
	for _, MESSAGE := range this.instructions {
		switch COMMAND := MESSAGE.(type) {
		case *commands.TrackOutcome:
			this.trackOutcome(COMMAND)
		case *commands.UpdateOutcomeTitle:
			this.updateOutcomeTitle(COMMAND)
		case *commands.UpdateOutcomeExplanation:
			this.updateOutcomeExplanation(COMMAND)
		case *commands.UpdateOutcomeDescription:
			this.updateOutcomeDescription(COMMAND)
		case *commands.DeleteOutcome:
			this.deleteOutcome(COMMAND)
		}
	}
}

func (this *Task) trackOutcome(_command *commands.TrackOutcome) {
	ID := this.nextID()
	AGGREGATE := this.aggregate(ID)
	_command.Result.ID = ID
	_command.Result.Error = AGGREGATE.TrackOutcome(ID, _command.Title)
}

func (this *Task) updateOutcomeTitle(_command *commands.UpdateOutcomeTitle) {
	AGGREGATE := this.aggregate(_command.OutcomeID)
	_command.Result.Error = AGGREGATE.UpdateOutcomeTitle(_command.UpdatedTitle)
}

func (this *Task) updateOutcomeExplanation(_command *commands.UpdateOutcomeExplanation) {
	AGGREGATE := this.aggregate(_command.OutcomeID)
	_command.Result.Error = AGGREGATE.UpdateOutcomeExplanation(_command.UpdatedExplanation)
}

func (this *Task) updateOutcomeDescription(_command *commands.UpdateOutcomeDescription) {
	AGGREGATE := this.aggregate(_command.OutcomeID)
	_command.Result.Error = AGGREGATE.UpdateOutcomeDescription(_command.UpdatedDescription)
}

func (this *Task) deleteOutcome(_command *commands.DeleteOutcome) {
	AGGREGATE := this.aggregate(_command.OutcomeID)
	_command.Result.Error = AGGREGATE.DeleteOutcome()
}

func (this *Task) publishResults() {
	for _, AGGREGATE := range this.aggregates {
		this.AddPendingWrites(AGGREGATE.TransferResults()...)
	}
}
