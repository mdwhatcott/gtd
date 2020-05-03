package eventstore

import (
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/util/errors"
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
		root, ok := event.(storage.Identifier)
		if !ok {
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(event)))
		}
		this.persist(root)
	}
}

func (this *Writer) persist(root storage.Identifier) {
	writer := this.writer(root)
	defer this.close(writer)
	err := this.encoder(writer).Encode(root)
	if err != nil {
		panic(errors.Wrap(ErrUnexpectedWriteError, err))
	}
}

func (this *Writer) close(writer io.WriteCloser) {
	err := writer.Close()
	if err != nil {
		panic(errors.Wrap(ErrUnexpectedCloseError, err))
	}
}
