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
	readErrs  map[string]error
	closeErrs map[string]error
	decodeErr error
}

func (this *ReaderFixture) readerFunc(id storage.Identifier) io.ReadCloser {
	return this.readers[id.ID()]
}

func (this *ReaderFixture) decoderFunc(reader io.Reader) storage.Decoder {
	return NewFakeDecoder(reader, this.decodeErr)
}

func (this *ReaderFixture) Setup() {
	this.readers = make(map[string]*FakeReader)
	this.readErrs = make(map[string]error)
	this.closeErrs = make(map[string]error)
	this.reader = NewReader(this.readerFunc, this.decoderFunc)
}

func (this *ReaderFixture) read(id string) (events []interface{}) {
	query := &storage.OutcomeEventStream{OutcomeID: id}
	this.reader.Read(query)
	return gather(query.Result.Stream)
}

func gather(stream chan interface{}) (all []interface{}) {
	for event := range stream {
		all = append(all, event)
	}
	return all
}

func (this *ReaderFixture) TestRead_UnrecognizedQueryType_PANIC() {
	action := func() { this.reader.Read(42) }
	result := recovered(action)
	this.So(result, should.BeError, "unrecognized query type: <int>")
}

func (this *ReaderFixture) TestRead() {
	this.readers["A"] = NewFakeReader("1\n2\n3")
	query := &storage.OutcomeEventStream{OutcomeID: "A"}

	this.reader.Read(query)

	this.So(gather(query.Result.Stream), should.Resemble, []interface{}{1, 2, 3})
	this.So(this.readers["A"].closed, should.Equal, 1)
}
