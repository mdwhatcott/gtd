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
func (this *Task) get(_id commands.Identifiable) *Aggregate {
	return this.getByID(_id.ID())
}
func (this *Task) getByID(_id string) *Aggregate {
	AGGREGATE, FOUND := this.aggregates[_id]
	if !FOUND {
		AGGREGATE = NewAggregate(this.clock.UTCNow(), this.log)
		this.aggregates[_id] = AGGREGATE
	}
	return AGGREGATE
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
	for ID, QUERY := range this.queries {
		this.getByID(ID).Replay(QUERY.Result.Events...)
	}
}
func (this *Task) processInstructions() {
	for _, MESSAGE := range this.instructions {
		switch COMMAND := MESSAGE.(type) {

		case *commands.TrackOutcome:
			ID := this.nextID()
			AGGREGATE := this.getByID(ID)
			COMMAND.Result.ID = ID
			COMMAND.Result.Error = AGGREGATE.TrackOutcome(ID, COMMAND.Title)

		case *commands.UpdateOutcomeTitle:
			COMMAND.Result.Error = this.get(COMMAND).UpdateOutcomeTitle(COMMAND.UpdatedTitle)

		case *commands.UpdateOutcomeExplanation:
			COMMAND.Result.Error = this.get(COMMAND).UpdateOutcomeExplanation(COMMAND.UpdatedExplanation)

		case *commands.UpdateOutcomeDescription:
			COMMAND.Result.Error = this.get(COMMAND).UpdateOutcomeDescription(COMMAND.UpdatedDescription)

		case *commands.DeleteOutcome:
			COMMAND.Result.Error = this.get(COMMAND).DeleteOutcome()

		case *commands.DeclareOutcomeRealized:
			COMMAND.Result.Error = this.get(COMMAND).DeclareOutcomeRealized()

		case *commands.DeclareOutcomeFixed:
			COMMAND.Result.Error = this.get(COMMAND).DeclareOutcomeFixed()
		}
	}
}
func (this *Task) publishResults() {
	for _, AGGREGATE := range this.aggregates {
		this.AddPendingWrites(AGGREGATE.TransferResults()...)
	}
}
