package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/core/events"
	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/util/fake"
)

func TestWriterFixture(t *testing.T) {
	gunit.Run(new(WriterFixture), t)
}

type WriterFixture struct {
	*gunit.Fixture

	writer    *Writer
	inner     *fake.Writer
	writeErr  error
	closeErr  error
	encodeErr error
}

func (this *WriterFixture) writerFunc() io.WriteCloser {
	if this.inner == nil {
		this.inner = fake.NewWriter(this.writeErr, this.closeErr)
	}
	return this.inner
}

func (this *WriterFixture) encoderFunc(_writer io.Writer) storage.Encoder {
	return fake.NewEncoder(_writer, this.encodeErr)
}

func (this *WriterFixture) Setup() {
	this.writer = NewWriter(this.encoderFunc, this.writerFunc)
}

func (this *WriterFixture) TestWrite_UnrecognizedEventType_PANIC() {
	ACTION := func() { this.writer.Write(42) }
	RESULT := recovered(ACTION)
	this.So(RESULT, should.Wrap, ErrUnrecognizedType)
}

func (this *WriterFixture) TestWrite_PersistsEncodedEventsToWriter() {
	this.writer.Write(outcomeTracked, outcomeUpdated)
	ACTUAL := this.inner.Lines()
	this.So(ACTUAL, should.Resemble, []string{"OutcomeTrackedV1", "OutcomeTitleUpdatedV1"})
	this.So(this.inner.CloseCount(), should.Equal, 2)
}

func (this *WriterFixture) TestWrite_ErrFromWriter_PANIC() {
	this.writeErr = errGophers
	ACTION := func() { this.writer.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

func (this *WriterFixture) TestWrite_ErrFromEncoder_PANIC() {
	this.encodeErr = errGophers
	ACTION := func() { this.writer.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

func (this *WriterFixture) TestWrite_ErrFromWriterClose_PANIC() {
	this.closeErr = errGophers
	ACTION := func() { this.writer.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedCloseError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

var (
	outcomeTracked = events.OutcomeTrackedV1{OutcomeID: "OutcomeID"}
	outcomeUpdated = events.OutcomeTitleUpdatedV1{OutcomeID: "OutcomeID"}
)
