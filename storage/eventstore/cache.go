package eventstore

import (
	"reflect"

	"github.com/smartystreets/joyride/v2"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/errors"
)

type Cache struct {
	log    *logging.Logger
	cached []interface{}
	writer joyride.StorageWriter
}

func NewCache(reader joyride.StorageReader, writer joyride.StorageWriter) *Cache {
	return &Cache{cached: warmUp(reader), writer: writer}
}

func warmUp(reader joyride.StorageReader) (cached_ []interface{}) {
	query := &storage.EventStream{}
	reader.Read(query)
	for event := range query.Result.Events {
		cached_ = append(cached_, event)
	}
	return cached_
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

func (this *Cache) stream(stream chan interface{}, filter string) {
	defer close(stream)

	COUNT := 0
	for _, VALUE := range this.cached {
		IDENTIFIABLE, OK := VALUE.(storage.Identifiable)
		if !OK {
			this.log.Println(errors.Wrap(ErrUnidentifiableType, reflect.TypeOf(VALUE)))
			break
		}

		if filter == "" || IDENTIFIABLE.ID() == filter {
			COUNT++
			stream <- VALUE
		}
	}
}

func (this *Cache) Write(i ...interface{}) {
	this.cached = append(this.cached, i...)
	this.writer.Write(i...)
}
