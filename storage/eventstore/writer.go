package eventstore

import (
	"io"
	"reflect"

	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/errors"
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
	for _, EVENT := range events {
		ROOT, OK := EVENT.(storage.Identifiable)
		if !OK {
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(EVENT)))
		}
		this.persist(ROOT)
	}
}

func (this *Writer) persist(root storage.Identifiable) {
	WRITER := this.writer()
	defer this.close(WRITER)

	ENCODER := this.encoder(WRITER)
	ERR := ENCODER.Encode(root)
	if ERR != nil {
		panic(errors.Wrap(ErrUnexpectedWriteError, ERR))
	}
}

func (this *Writer) close(writer io.WriteCloser) {
	ERR := writer.Close()
	if ERR != nil {
		panic(errors.Wrap(ErrUnexpectedCloseError, ERR))
	}
}
