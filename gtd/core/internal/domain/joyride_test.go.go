package domain

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/storage/queries"
)

type FakeShell struct {
	*gunit.Fixture
	writes   []interface{}
	messages []interface{}
	reads    map[string][]interface{} // map[user-id][]event
}

func NewFakeShell(fixture *gunit.Fixture) *FakeShell {
	return &FakeShell{
		Fixture: fixture,
		reads:   make(map[string][]interface{}),
	}
}

func (this *FakeShell) Dispatch(values ...interface{}) {
	this.messages = append(this.messages, values...)
}

func (this *FakeShell) Write(values ...interface{}) {
	this.writes = append(this.writes, values...)
}

func (this *FakeShell) Read(values ...interface{}) {
	for _, value := range values {
		switch query := value.(type) {
		case *queries.EventStream:
			query.Result.Stream = load(this.reads[query.UserID])
		}
	}
}
func load(events []interface{}) chan interface{} {
	stream := make(chan interface{})
	go populate(events, stream)
	return stream
}
func populate(events []interface{}, stream chan interface{}) {
	for _, event := range events {
		stream <- event
	}
}

func (this *FakeShell) PrepareReadResults(user string, results ...interface{}) {
	this.reads[user] = append(this.reads[user], results)
}

func (this *FakeShell) AssertOutput(expected ...interface{}) {
	if !this.So(len(this.writes), should.Equal, len(expected)) {
		return
	}
	var failed bool
	for _, e := range expected {
		failed = failed || !this.So(this.writes, should.Contain, e)
	}
	if failed {
		return
	}
	this.So(this.writes, should.Resemble, this.messages)
}
