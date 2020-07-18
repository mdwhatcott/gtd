package eventstore

import (
	"context"
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core/events"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/util/fake"
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

func (this *WriterFixture) encoderFunc(writer io.Writer) storage.Encoder {
	return fake.NewEncoder(writer, this.encodeErr)
}

func (this *WriterFixture) Setup() {
	this.writer = NewWriter(this.encoderFunc, this.writerFunc)
}

func (this *WriterFixture) write(v ...interface{}) {
	this.writer.Write(context.Background(), v...)
}

func (this *WriterFixture) TestWrite_UnrecognizedEventType_PANIC() {
	ACTION := func() { this.write(42) }
	RESULT := recovered(ACTION)
	this.So(RESULT, should.Wrap, ErrUnrecognizedType)
}

func (this *WriterFixture) TestWrite_PersistsEncodedEventsToWriter() {
	this.write(outcomeTracked, outcomeUpdated)
	ACTUAL := this.inner.Lines()
	this.So(ACTUAL, should.Resemble, []string{"OutcomeTrackedV1", "OutcomeTitleUpdatedV1"})
	this.So(this.inner.CloseCount(), should.Equal, 2)
}

func (this *WriterFixture) TestWrite_ErrFromWriter_PANIC() {
	this.writeErr = errGophers
	ACTION := func() { this.write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

func (this *WriterFixture) TestWrite_ErrFromEncoder_PANIC() {
	this.encodeErr = errGophers
	ACTION := func() { this.write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

func (this *WriterFixture) TestWrite_ErrFromWriterClose_PANIC() {
	this.closeErr = errGophers
	ACTION := func() { this.write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedCloseError)
	this.So(this.inner.CloseCount(), should.Equal, 1)
}

var (
	outcomeTracked = events.OutcomeTrackedV1{OutcomeID: "OutcomeID"}
	outcomeUpdated = events.OutcomeTitleUpdatedV1{OutcomeID: "OutcomeID"}
)
