package projections

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type Outcome struct {
	Updated     time.Time
	ID          string
	Title       string
	Status      string
	Explanation string
	Description string
}

func NewOutcome() *Outcome {
	return &Outcome{}
}

func (this *Outcome) Apply(_messages ...interface{}) {
	for _, MESSAGE := range _messages {
		switch EVENT := MESSAGE.(type) {
		case events.OutcomeTrackedV1:
			this.Updated = EVENT.Timestamp
			this.ID = EVENT.OutcomeID
			this.Title = EVENT.Title
		case events.OutcomeFixedV1:
			this.Status = "fixed"
		}
	}
}
