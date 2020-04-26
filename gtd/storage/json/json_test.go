package json

import (
	"bytes"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/util/date"
)

func TestEncoderFixture(t *testing.T) {
	gunit.Run(new(EncoderFixture), t)
}

type EncoderFixture struct {
	*gunit.Fixture

	writer  *bytes.Buffer
	encoder *Encoder
}

func (this *EncoderFixture) Setup() {
	this.writer = new(bytes.Buffer)
	this.encoder = NewEncoder(this.writer)
}

func (this *EncoderFixture) TestEventSerialization() {
	event := events.OutcomeTrackedV1{
		Timestamp: date.YMD(2020, 1, 1),
		OutcomeID: "OutcomeID",
		Title:     "Title",
	}
	err := this.encoder.Encode(event)

	this.So(err, should.BeNil)
	this.So("\n"+this.writer.String(), should.Equal, `
"events.OutcomeTrackedV1"
{
  "timestamp": "2020-01-01T00:00:00Z",
  "outcome_id": "OutcomeID",
  "title": "Title"
}
`)
}

func (this *EncoderFixture) TestErr() {
	err := this.encoder.Encode(make(chan int))
	this.So(err, should.NotBeNil)
	this.So(this.writer.String(), should.Equal, `"chan int"`+"\n")
}
