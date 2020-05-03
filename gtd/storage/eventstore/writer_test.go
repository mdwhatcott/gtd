package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/storage"
)

func TestWriterFixture(t *testing.T) {
	gunit.Run(new(WriterFixture), t)
}

type WriterFixture struct {
	*gunit.Fixture

	repo      *Writer
	history   map[string][]interface{}
	writers   map[string]*FakeWriter
	writeErrs map[string]error
	closeErrs map[string]error
	encodeErr error
}

func (this *WriterFixture) writerFunc(root storage.Identifier) io.WriteCloser {
	id := root.ID()
	writer, found := this.writers[id]
	if !found {
		writer = NewFakeWriter(this.writeErrs[id], this.closeErrs[id])
		this.writers[id] = writer
	}
	return writer
}

func (this *WriterFixture) encoderFunc(writer io.Writer) storage.Encoder {
	return NewFakeEncoder(writer, this.encodeErr)
}

func (this *WriterFixture) Setup() {
	this.history = make(map[string][]interface{})
	this.writers = make(map[string]*FakeWriter)
	this.writeErrs = make(map[string]error)
	this.closeErrs = make(map[string]error)
	this.repo = NewWriter(this.encoderFunc, this.writerFunc)
}

func (this *WriterFixture) TestWrite_UnrecognizedEventType_PANIC() {
	action := func() { this.repo.Write(42) }
	result := recovered(action)
	this.So(result, should.Wrap, ErrUnrecognizedType)
}

func (this *WriterFixture) TestWrite_PersistsEncodedEventsToWriter() {
	this.repo.Write(outcomeTracked, outcomeUpdated)
	actual := this.writers["OutcomeID"].Lines()
	this.So(actual, should.Resemble, []string{"OutcomeTrackedV1", "OutcomeTitleUpdatedV1"})
}

func (this *WriterFixture) TestWrite_ErrFromWriter_PANIC() {
	this.writeErrs["OutcomeID"] = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, ErrUnexpectedWriteError)
}

func (this *WriterFixture) TestWrite_ErrFromEncoder_PANIC() {
	this.encodeErr = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, ErrUnexpectedWriteError)
}

func (this *WriterFixture) TestWrite_ErrFromWriterClose_PANIC() {
	this.closeErrs["OutcomeID"] = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, ErrUnexpectedCloseError)
}

var (
	outcomeTracked = events.OutcomeTrackedV1{OutcomeID: "OutcomeID"}
	outcomeUpdated = events.OutcomeTitleUpdatedV1{OutcomeID: "OutcomeID"}
)
