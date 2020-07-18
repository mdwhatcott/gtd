package fake

import (
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/v3/storage"
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

func (this *Joyride) Write(values ...interface{}) {
	this.Writes = append(this.Writes, values...)
}

func (this *Joyride) Read(values ...interface{}) {
	for _, VALUE := range values {
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

func (this *Joyride) PrepareReadResults(id string, results ...interface{}) {
	this.reads[id] = append(this.reads[id], results...)
}
