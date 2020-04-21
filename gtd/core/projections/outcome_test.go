package projections

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
)

func TestOutcomeFixture(t *testing.T) {
	gunit.Run(new(OutcomeFixture), t)
}

type OutcomeFixture struct {
	*gunit.Fixture
	projection *Outcome
}

func (this *OutcomeFixture) Setup() {
	this.projection = NewOutcome()
}

func date(ymd ...int) time.Time {
	return time.Date(ymd[0], time.Month(ymd[1]), ymd[2], 0, 0, 0, 0, time.UTC)
}

func (this *OutcomeFixture) TestBlankWhenFirstInstantiated() {
	this.So(this.projection, should.Resemble, &Outcome{
		Updated:     time.Time{},
		ID:          "",
		Title:       "",
		Explanation: "",
		Description: "",
	})
}

func (this *OutcomeFixture) TestOutcomeTracked() {
	this.projection.Apply(events.OutcomeTrackedV1{
		Timestamp: date(2020, 1, 1),
		OutcomeID: "id-1",
		Title:     "title",
	})

	this.So(this.projection, should.Resemble, &Outcome{
		Updated:     date(2020, 1, 1),
		ID:          "id-1",
		Title:       "title",
		Explanation: "",
		Description: "",
	})
}
