package eventstore

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type ReadWriter struct{ Dependencies }

type Dependencies struct {
	encoder storage.EncoderFunc
	writer  storage.WriterFunc
	history map[string][]interface{}
}

func NewReadWriter(dependencies Dependencies) *ReadWriter {
	return &ReadWriter{Dependencies: dependencies}
}
func (this *ReadWriter) Read(queries ...interface{}) {
	for _, query := range queries {
		switch query := query.(type) {
		case *storage.OutcomeEventStream:
			query.Result.Events = this.history[query.OutcomeID]
		default:
			panic(fmt.Errorf("unrecognized query type: <%v>", reflect.TypeOf(query)))
		}
	}
}
func (this *ReadWriter) Write(events ...interface{}) {
	for _, event := range events {
		root, ok := event.(storage.AggregateRoot)
		if !ok {
			panic(fmt.Errorf("unrecognized event type: <%v>", reflect.TypeOf(event)))
		}
		this.persist(root)
		this.append(root)
	}
}
func (this *ReadWriter) append(event storage.AggregateRoot) {
	id := event.ID()
	events := this.history[id]
	events = append(events, event)
	this.history[id] = events
}
func (this *ReadWriter) persist(root storage.AggregateRoot) {
	writer := this.writer(root)
	defer this.close(writer)
	err := this.encoder(writer).Encode(root)
	if err != nil {
		panic(fmt.Errorf("persistence error: %w", err))
	}
}
func (this *ReadWriter) close(writer io.WriteCloser) {
	err := writer.Close()
	if err != nil {
		panic(fmt.Errorf("persistence error (on close): %w", err))
	}
}
