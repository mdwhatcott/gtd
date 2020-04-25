package eventstore

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type OutcomeRepository struct{ OutcomeRepositoryState }

type OutcomeRepositoryState struct {
	encoder storage.EncoderFunc
	writer  storage.WriterFunc
	history map[string][]interface{}
}

func NewOutcomeRepository(state OutcomeRepositoryState) *OutcomeRepository {
	return &OutcomeRepository{OutcomeRepositoryState: state}
}
func (this *OutcomeRepository) Read(queries ...interface{}) {
	for _, query := range queries {
		switch query := query.(type) {
		case *storage.OutcomeEventStream:
			query.Result.Events = this.history[query.OutcomeID]
		default:
			panic(fmt.Errorf("unrecognized query type: <%v>", reflect.TypeOf(query)))
		}
	}
}
func (this *OutcomeRepository) Write(_events ...interface{}) {
	for _, event := range _events {
		root, ok := event.(storage.AggregateRoot)
		if !ok {
			panic(fmt.Errorf("unrecognized event type: <%v>", reflect.TypeOf(event)))
		}
		this.persist(root)
		this.append(root)
	}
}
func (this *OutcomeRepository) append(event storage.AggregateRoot) {
	id := event.ID()
	events := this.history[id]
	events = append(events, event)
	this.history[id] = events
}
func (this *OutcomeRepository) persist(root storage.AggregateRoot) {
	writer := this.writer(root)
	defer this.close(writer)
	err := this.encoder(writer).Encode(root)
	if err != nil {
		panic(fmt.Errorf("persistence error: %w", err))
	}
}
func (this *OutcomeRepository) close(writer io.WriteCloser) {
	err := writer.Close()
	if err != nil {
		panic(fmt.Errorf("persistence error (on close): %w", err))
	}
}
