package eventstore

import (
	"io"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

func TestReaderFixture(t *testing.T) {
	gunit.Run(new(ReaderFixture), t)
}

type ReaderFixture struct {
	*gunit.Fixture
	reader    *Reader
	readers   map[string]*FakeReader
	decodeErr error
}

func (this *ReaderFixture) readerFunc(_id storage.Identifier) io.ReadCloser {
	return this.readers[_id.ID()]
}

func (this *ReaderFixture) decoderFunc(_reader io.Reader) storage.Decoder {
	return NewFakeDecoder(_reader, this.decodeErr)
}

func (this *ReaderFixture) Setup() {
	this.readers = make(map[string]*FakeReader)
	this.reader = NewReader(this.readerFunc, this.decoderFunc)
}

func (this *ReaderFixture) read(_id string) (events_ []interface{}) {
	query := &storage.OutcomeEventStream{OutcomeID: _id}
	this.reader.Read(query)
	return gather(query.Result.Stream)
}

func gather(_stream chan interface{}) (all_ []interface{}) {
	for EVENT := range _stream {
		all_ = append(all_, EVENT)
	}
	return all_
}

func (this *ReaderFixture) TestRead_UnrecognizedQueryType_PANIC() {
	ACTION := func() { this.reader.Read(42) }
	RESULT := recovered(ACTION)
	this.So(RESULT, should.Wrap, ErrUnrecognizedType)
}

func (this *ReaderFixture) TestRead() {
	this.readers["A"] = NewFakeReader("1\n2\n3")
	QUERY := &storage.OutcomeEventStream{OutcomeID: "A"}

	this.reader.Read(QUERY)

	this.So(gather(QUERY.Result.Stream), should.Resemble, []interface{}{1, 2, 3})
	this.So(this.readers["A"].closed, should.Equal, 1)
}

func (this *ReaderFixture) TestReadErr() {
	this.readers["A"] = NewFakeReader("1\n2\n3")
	this.readers["A"].readErr = errGophers
	QUERY := &storage.OutcomeEventStream{OutcomeID: "A"}

	this.reader.Read(QUERY)

	RESULTS := gather(QUERY.Result.Stream)
	if this.So(RESULTS, should.HaveLength, 1) {
		this.So(RESULTS[0], should.Wrap, ErrUnexpectedReadError)
	}
}

func (this *ReaderFixture) TestCloseErr() {
	this.readers["A"] = NewFakeReader("1\n2\n3")
	this.readers["A"].closeErr = errGophers
	QUERY := &storage.OutcomeEventStream{OutcomeID: "A"}

	this.reader.Read(QUERY)

	RESULTS := gather(QUERY.Result.Stream)
	if this.So(RESULTS, should.HaveLength, 4) {
		this.So(RESULTS[:3], should.Resemble, []interface{}{1, 2, 3})
		this.So(RESULTS[3], should.Wrap, ErrUnexpectedCloseError)
	}
}
