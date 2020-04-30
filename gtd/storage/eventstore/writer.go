package eventstore

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Writer struct {
	encoder storage.EncoderFunc
	writer  storage.WriterFunc
}

func NewWriter(encoder storage.EncoderFunc, writer storage.WriterFunc) *Writer {
	return &Writer{
		encoder: encoder,
		writer:  writer,
	}
}

func (this *Writer) Write(events ...interface{}) {
	for _, event := range events {
		root, ok := event.(storage.AggregateRoot)
		if !ok {
			panic(fmt.Errorf("unrecognized event type: <%v>", reflect.TypeOf(event)))
		}
		this.persist(root)
	}
}

func (this *Writer) persist(root storage.AggregateRoot) {
	writer := this.writer(root)
	defer this.close(writer)
	err := this.encoder(writer).Encode(root)
	if err != nil {
		panic(fmt.Errorf("persistence error: %w", err))
	}
}

func (this *Writer) close(writer io.WriteCloser) {
	err := writer.Close()
	if err != nil {
		panic(fmt.Errorf("persistence error (on close): %w", err))
	}
}
