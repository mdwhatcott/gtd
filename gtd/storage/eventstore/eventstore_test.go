package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/storage"
)

func TestOutcomeRepositoryFixture(t *testing.T) {
	gunit.Run(new(OutcomeRepositoryFixture), t)
}

type OutcomeRepositoryFixture struct {
	*gunit.Fixture

	repo      *ReadWriter
	history   map[string][]interface{}
	writers   map[string]*FakeWriter
	writeErrs map[string]error
	closeErrs map[string]error
	encodeErr error
}

func (this *OutcomeRepositoryFixture) writerFunc(root storage.AggregateRoot) io.WriteCloser {
	id := root.ID()
	writer, found := this.writers[id]
	if !found {
		writer = NewFakeWriter(this.writeErrs[id], this.closeErrs[id])
		this.writers[id] = writer
	}
	return writer
}
func (this *OutcomeRepositoryFixture) encoderFunc(writer io.Writer) storage.Encoder {
	return NewFakeEncoder(writer, this.encodeErr)
}
func (this *OutcomeRepositoryFixture) read(id string) (events []interface{}) {
	query := &storage.OutcomeEventStream{OutcomeID: id}
	this.repo.Read(query)
	return query.Result.Events
}
func (this *OutcomeRepositoryFixture) Setup() {
	this.history = make(map[string][]interface{})
	this.writers = make(map[string]*FakeWriter)
	this.writeErrs = make(map[string]error)
	this.closeErrs = make(map[string]error)
	this.repo = NewReadWriter(Dependencies{
		encoder: this.encoderFunc,
		writer:  this.writerFunc,
		history: this.history,
	})
}
func (this *OutcomeRepositoryFixture) TestRead_UnrecognizedQueryType_PANIC() {
	action := func() { this.repo.Read(42) }
	result := recovered(action)
	this.So(result, should.BeError, "unrecognized query type: <int>")
}
func (this *OutcomeRepositoryFixture) TestRead_OutcomeEventStream_ResultPopulated() {
	this.history["a"] = []interface{}{1, 2, 3}
	stream := this.read("a")
	this.So(stream, should.Resemble, []interface{}{1, 2, 3})
}
func (this *OutcomeRepositoryFixture) TestRead_OutcomeEventStream_NoHistory_ResultNil() {
	this.history["a"] = nil
	steam := this.read("a")
	this.So(steam, should.BeNil)
}
func (this *OutcomeRepositoryFixture) TestWrite_UnrecognizedEventType_PANIC() {
	action := func() { this.repo.Write(42) }
	result := recovered(action)
	this.So(result, should.BeError, "unrecognized event type: <int>")
}
func (this *OutcomeRepositoryFixture) TestWrite_AppendsToInMemoryCollection() {
	this.repo.Write(outcomeTracked, outcomeUpdated)
	actual := this.read(outcomeTracked.OutcomeID)
	this.So(actual, should.Resemble, []interface{}{outcomeTracked, outcomeUpdated})
}
func (this *OutcomeRepositoryFixture) TestWrite_PersistsEncodedEventsToWriter() {
	this.repo.Write(outcomeTracked, outcomeUpdated)
	actual := this.writers["OutcomeID"].Lines()
	this.So(actual, should.Resemble, []string{"OutcomeTrackedV1", "OutcomeTitleUpdatedV1"})
}
func (this *OutcomeRepositoryFixture) TestWrite_ErrFromWriter_PANIC() {
	this.writeErrs["OutcomeID"] = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, errGophers)
}
func (this *OutcomeRepositoryFixture) TestWrite_ErrFromEncoder_PANIC() {
	this.encodeErr = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, errGophers)
}
func (this *OutcomeRepositoryFixture) TestWrite_ErrFromWriterClose_PANIC() {
	this.closeErrs["OutcomeID"] = errGophers
	action := func() { this.repo.Write(outcomeTracked) }
	this.So(recovered(action), should.Wrap, errGophers)
}

var (
	outcomeTracked = events.OutcomeTrackedV1{OutcomeID: "OutcomeID"}
	outcomeUpdated = events.OutcomeTitleUpdatedV1{OutcomeID: "OutcomeID"}
)
