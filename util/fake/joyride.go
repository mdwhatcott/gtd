package fake

import (
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/storage"
)

type Joyride struct {
	log    *logging.Logger
	Writes []interface{}
	reads  map[string][]interface{}
}

func NewJoyride(log *logging.Logger) *Joyride {
	return &Joyride{
		log:   log,
		reads: make(map[string][]interface{}),
	}
}

func (this *Joyride) Write(_values ...interface{}) {
	this.Writes = append(this.Writes, _values...)
}

func (this *Joyride) Read(_values ...interface{}) {
	for _, VALUE := range _values {
		this.log.Println("Reading value:", VALUE)
		switch QUERY := VALUE.(type) {
		case *storage.OutcomeEventStream:
			this.log.Println("Reading outcome event stream...", QUERY.OutcomeID)
			QUERY.Result.Events = this.reads[QUERY.OutcomeID]
		}
	}
}

func (this *Joyride) PrepareReadResults(_id string, _results ...interface{}) {
	this.reads[_id] = append(this.reads[_id], _results...)
	this.log.Println("Read:", _id, _results)
}
