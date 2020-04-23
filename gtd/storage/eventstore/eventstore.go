package eventstore

import (
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type OutcomeRepository struct {
	log *logging.Logger

	encoder storage.NewEncoder
	writer  storage.Writer
	roots   map[string][]interface{}
}

func NewOutcomeRepository(encoder storage.NewEncoder, writer storage.Writer) *OutcomeRepository {
	return &OutcomeRepository{
		encoder: encoder,
		writer:  writer,
		roots:   make(map[string][]interface{}),
	}
}
func (this *OutcomeRepository) Read(queries ...interface{}) {
	for _, query := range queries {
		switch query := query.(type) {
		case storage.OutcomeEventStream:
			query.Result.Events = this.roots[query.OutcomeID]
		default:
			this.log.Panicf("UNRECOGNIZED QUERY TYPE") // TODO: better log
		}
	}
}
func (this *OutcomeRepository) Write(_events ...interface{}) {
	for _, event := range _events {
		root, ok := event.(storage.AggregateRoot)
		if !ok {
			this.log.Panicf("UNRECOGNIZED EVENT TYPE") // TODO: better log
		}
		this.persist(root)
		this.append(root)
	}
}
func (this *OutcomeRepository) append(event storage.AggregateRoot) {
	id := event.ID()
	events := this.roots[id]
	events = append(events, event)
	this.roots[id] = events
}
func (this *OutcomeRepository) persist(root storage.AggregateRoot) {
	encoder := this.encoder(this.writer(root))
	err := encoder.Encode(root)
	if err != nil {
		this.log.Panicf("%s", err) // TODO: better log
	}
}
