package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/util/fake"
)

func TestWriterFixture(t *testing.T) {
	gunit.Run(new(WriterFixture), t)
}

type WriterFixture struct {
	*gunit.Fixture

	repo      *Writer
	history   map[string][]interface{}
	writers   map[string]*fake.Writer
	writeErrs map[string]error
	closeErrs map[string]error
	encodeErr error
}

func (this *WriterFixture) writerFunc(_root storage.Identifier) io.WriteCloser {
	ID := _root.ID()
	WRITER, FOUND := this.writers[ID]
	if !FOUND {
		WRITER = fake.NewWriter(this.writeErrs[ID], this.closeErrs[ID])
		this.writers[ID] = WRITER
	}
	return WRITER
}

func (this *WriterFixture) encoderFunc(_writer io.Writer) storage.Encoder {
	return fake.NewEncoder(_writer, this.encodeErr)
}

func (this *WriterFixture) Setup() {
	this.history = make(map[string][]interface{})
	this.writers = make(map[string]*fake.Writer)
	this.writeErrs = make(map[string]error)
	this.closeErrs = make(map[string]error)
	this.repo = NewWriter(this.encoderFunc, this.writerFunc)
}

func (this *WriterFixture) TestWrite_UnrecognizedEventType_PANIC() {
	ACTION := func() { this.repo.Write(42) }
	RESULT := recovered(ACTION)
	this.So(RESULT, should.Wrap, ErrUnrecognizedType)
}

func (this *WriterFixture) TestWrite_PersistsEncodedEventsToWriter() {
	this.repo.Write(outcomeTracked, outcomeUpdated)
	ACTUAL := this.writers["OutcomeID"].Lines()
	this.So(ACTUAL, should.Resemble, []string{"OutcomeTrackedV1", "OutcomeTitleUpdatedV1"})
}

func (this *WriterFixture) TestWrite_ErrFromWriter_PANIC() {
	this.writeErrs["OutcomeID"] = errGophers
	ACTION := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
}

func (this *WriterFixture) TestWrite_ErrFromEncoder_PANIC() {
	this.encodeErr = errGophers
	ACTION := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedWriteError)
}

func (this *WriterFixture) TestWrite_ErrFromWriterClose_PANIC() {
	this.closeErrs["OutcomeID"] = errGophers
	ACTION := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(ACTION), should.Wrap, ErrUnexpectedCloseError)
}

var (
	outcomeTracked = events.OutcomeTrackedV1{OutcomeID: "OutcomeID"}
	outcomeUpdated = events.OutcomeTitleUpdatedV1{OutcomeID: "OutcomeID"}
)
