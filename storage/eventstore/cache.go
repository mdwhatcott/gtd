package eventstore

import (
	"reflect"

	"github.com/smartystreets/joyride/v2"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/util/errors"
)

type Cache struct {
	log    *logging.Logger
	cached []interface{}
	writer joyride.StorageWriter
}

func NewCache(reader joyride.StorageReader, writer joyride.StorageWriter) *Cache {
	return &Cache{cached: warmUp(reader), writer: writer}
}

func warmUp(reader joyride.StorageReader) (cached []interface{}) {
	query := &storage.EventStream{}
	reader.Read(query)
	for event := range query.Result.Events {
		cached = append(cached, event)
	}
	return cached
}

func (this *Cache) Read(v ...interface{}) {
	for _, QUERY := range v {
		switch QUERY := QUERY.(type) {
		case *storage.EventStream:
			QUERY.Result.Events = make(chan interface{})
			go this.stream(QUERY.Result.Events, "")
		case *storage.OutcomeEventStream:
			QUERY.Result.Events = make(chan interface{})
			go this.stream(QUERY.Result.Events, QUERY.OutcomeID)
		default:
			panic(errors.Wrap(ErrUnrecognizedType, reflect.TypeOf(QUERY)))
		}
	}
}

func (this *Cache) stream(_stream chan interface{}, _filter string) {
	defer close(_stream)

	COUNT := 0
	for _, VALUE := range this.cached {
		IDENTIFIABLE, OK := VALUE.(storage.Identifiable)
		if !OK {
			this.log.Println(errors.Wrap(ErrUnidentifiableType, reflect.TypeOf(VALUE)))
			break
		}

		if _filter == "" || IDENTIFIABLE.ID() == _filter {
			COUNT++
			_stream <- VALUE
		}
	}
}

func (this *Cache) Write(i ...interface{}) {
	this.cached = append(this.cached, i...)
	this.writer.Write(i...)
}
