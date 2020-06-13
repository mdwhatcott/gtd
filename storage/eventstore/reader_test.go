package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/util/fake"
)

func TestReaderFixture(t *testing.T) {
	gunit.Run(new(ReaderFixture), t)
}

type ReaderFixture struct {
	*gunit.Fixture
	reader    *Reader
	inner     *fake.Reader
	decodeErr error
}

func (this *ReaderFixture) readerFunc() io.ReadCloser {
	return this.inner
}

func (this *ReaderFixture) decoderFunc(_reader io.Reader) storage.Decoder {
	return fake.NewDecoder(_reader, this.decodeErr)
}

func (this *ReaderFixture) Setup() {
	this.inner = fake.NewReader("1\n2\n3")
	this.reader = NewReader(this.readerFunc, this.decoderFunc)
	this.reader.log = logging.Capture()
}

func (this *ReaderFixture) TestRead_UnrecognizedQueryType_PANIC() {
	ACTION := func() { this.reader.Read(42) }
	RESULT := recovered(ACTION)
	this.So(RESULT, should.Wrap, ErrUnrecognizedType)
}

func (this *ReaderFixture) TestRead_OutcomeEventStream_EventsFilteredByID() {
	this.inner = fake.NewReader("1\n2\n2")
	QUERY := &storage.OutcomeEventStream{OutcomeID: "2"}

	this.reader.Read(QUERY)

	this.So(stream(QUERY.Result.Events), should.Resemble, []interface{}{
		fake.NewIdentifiable(2),
		fake.NewIdentifiable(2),
	})
	this.So(this.inner.Closed, should.Equal, 1)
}

func (this *ReaderFixture) TestRead_OutcomeEventStream_NonIdentifiableValueEncountered_PANIC() {
	this.inner = fake.NewReader("1\n2\n-1000") // The fake decoder treats negative numbers as unidentifiable.
	QUERY := &storage.OutcomeEventStream{OutcomeID: "2"}

	this.reader.Read(QUERY)

	this.So(stream(QUERY.Result.Events), should.Resemble, []interface{}{
		fake.NewIdentifiable(2),
	})
	this.So(this.reader.log.Log.String(), should.ContainSubstring, ErrUnidentifiableType.Error())
	this.So(this.inner.Closed, should.Equal, 1)
}

func (this *ReaderFixture) TestRead_OutcomeEventStreamErr() {
	this.inner.ReadErr = errGophers
	QUERY := new(storage.OutcomeEventStream)

	this.reader.Read(QUERY)

	this.So(stream(QUERY.Result.Events), should.BeEmpty)
	this.So(this.reader.log.Log.String(), should.ContainSubstring, ErrUnexpectedReadError.Error())
}

func (this *ReaderFixture) TestCloseErr() {
	this.inner.CloseErr = errGophers
	QUERY := new(storage.OutcomeEventStream)

	this.reader.Read(QUERY)

	this.So(stream(QUERY.Result.Events), should.Resemble, []interface{}{
		fake.NewIdentifiable(1),
		fake.NewIdentifiable(2),
		fake.NewIdentifiable(3),
	})
	this.So(this.reader.log.Log.String(), should.ContainSubstring, ErrUnexpectedCloseError.Error())
}

func (this *ReaderFixture) TestRead_EventStream_AllEventsIncluded() {
	QUERY := new(storage.EventStream)

	this.reader.Read(QUERY)

	this.So(stream(QUERY.Result.Events), should.Resemble, []interface{}{
		fake.NewIdentifiable(1),
		fake.NewIdentifiable(2),
		fake.NewIdentifiable(3),
	})
}

func stream(events chan interface{}) (streamed []interface{}) {
	for event := range events {
		streamed = append(streamed, event)
	}
	return streamed
}
