package projections

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/core/events"
)

type OutcomeDetail struct {
	Updated     time.Time
	ID          string
	Title       string
	Status      string
	Explanation string
	Description string
}

func NewOutcomeDetail() *OutcomeDetail {
	return &OutcomeDetail{}
}

func (this *OutcomeDetail) Apply(_messages ...interface{}) {
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
