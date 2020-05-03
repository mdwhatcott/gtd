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

func (this *Outcome) Apply(messages ...interface{}) {
	// TODO: panic if the event is an error value because something went very wrong in the load operation.
	for _, message := range messages {
		switch event := message.(type) {
		case events.OutcomeTrackedV1:
			this.Updated = event.Timestamp
			this.ID = event.OutcomeID
			this.Title = event.Title
		case events.OutcomeFixedV1:
			this.Status = "fixed"
		}
	}
}
