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
		switch QUERY := VALUE.(type) {
		case *storage.OutcomeEventStream:
			QUERY.Result.Events = make(chan interface{}, len(this.reads[QUERY.OutcomeID]))
			for _, EVENT := range this.reads[QUERY.OutcomeID] {
				QUERY.Result.Events <- EVENT
			}
			close(QUERY.Result.Events)
		}
	}
}

func (this *Joyride) PrepareReadResults(_id string, _results ...interface{}) {
	this.reads[_id] = append(this.reads[_id], _results...)
}
