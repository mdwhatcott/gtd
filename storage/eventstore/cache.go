package eventstore

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/errors"
)

type Cache struct {
	log    core.Logger
	cached []interface{}
	writer joyride.StorageWriter
}

func NewCache(logger core.Logger, reader joyride.StorageReader, writer joyride.StorageWriter) *Cache {
	return &Cache{
		log:    logger,
		cached: warmUp(reader),
		writer: writer,
	}
}

func warmUp(reader joyride.StorageReader) (cached_ []interface{}) {
	start := time.Now()
	query := &storage.EventStream{}
	reader.Read(context.Background(), query)
	for event := range query.Result.Events {
		cached_ = append(cached_, event)
	}
	log.Println()
	log.Printf("Cached warmed (%d events in %v)", len(cached_), time.Since(start))
	return cached_
}

func (this *Cache) Read(_ context.Context, v ...interface{}) {
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

func (this *Cache) Write(ctx context.Context, i ...interface{}) {
	this.cached = append(this.cached, i...)
	this.writer.Write(ctx, i...)
}
