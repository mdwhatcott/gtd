package projections

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/util/date"
)

func TestOutcomeFixture(t *testing.T) {
	gunit.Run(new(OutcomeFixture), t)
}

type OutcomeFixture struct {
	*gunit.Fixture
	projection *OutcomeDetail
}

func (this *OutcomeFixture) Setup() {
	this.projection = NewOutcomeDetail()
}

func (this *OutcomeFixture) TestBlankWhenFirstInstantiated() {
	this.So(this.projection, should.Resemble, &OutcomeDetail{
		Updated:     time.Time{},
		ID:          "",
		Title:       "",
		Explanation: "",
		Description: "",
	})
}

func (this *OutcomeFixture) TestOutcomeTracked() {
	this.projection.Apply(events.OutcomeTrackedV1{
		Timestamp: date.YMD(2020, 1, 1),
		OutcomeID: "id-1",
		Title:     "title",
	})

	this.So(this.projection, should.Resemble, &OutcomeDetail{
		Updated:     date.YMD(2020, 1, 1),
		ID:          "id-1",
		Title:       "title",
		Explanation: "",
		Description: "",
	})
}

func (this *OutcomeFixture) TestOutcomeFixed() {
	this.projection.Apply(
		events.OutcomeTrackedV1{
			Timestamp: date.YMD(2020, 1, 1),
			OutcomeID: "id-1",
			Title:     "title",
		},
		events.OutcomeFixedV1{
			Timestamp: date.YMD(2020, 1, 1),
			OutcomeID: "id-1",
		},
	)

	this.So(this.projection, should.Resemble, &OutcomeDetail{
		Updated:     date.YMD(2020, 1, 1),
		ID:          "id-1",
		Title:       "title",
		Status:      "fixed",
		Explanation: "",
		Description: "",
	})
}
