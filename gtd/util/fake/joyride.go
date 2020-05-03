package fake

import (
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Joyride struct {
	log    *logging.Logger
	Writes []interface{}
	reads  map[string][]interface{}
}

func NewJoyride() *Joyride {
	return &Joyride{reads: make(map[string][]interface{})}
}

func (this *Joyride) Write(_values ...interface{}) {
	this.Writes = append(this.Writes, _values...)
}

func (this *Joyride) Read(_values ...interface{}) {
	this.log.Println("Reading:", _values)
	for _, VALUE := range _values {
		this.log.Println("Reading value:", VALUE)
		switch QUERY := VALUE.(type) {
		case *storage.OutcomeEventStream:
			this.log.Println("Reading outcome event stream...", QUERY.OutcomeID)
			QUERY.Result.Stream = make(chan interface{})
			go load(QUERY.Result.Stream, this.reads[QUERY.OutcomeID])
		}
	}
}

func load(_stream chan interface{}, _events []interface{}) {
	defer close(_stream)
	for _, EVENT := range _events {
		_stream <- EVENT
	}
}

func (this *Joyride) PrepareReadResults(_id string, _results ...interface{}) {
	this.reads[_id] = append(this.reads[_id], _results...)
	this.log.Println("Read:", _id, _results)
}
