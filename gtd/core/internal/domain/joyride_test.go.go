package domain

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/storage/queries"
)

type FakeShell struct {
	*gunit.Fixture

	writes   []interface{}
	reads    []interface{}
}

func NewFakeShell(fixture *gunit.Fixture) *FakeShell {
	return &FakeShell{Fixture: fixture}
}

func (this *FakeShell) Write(values ...interface{}) {
	this.writes = append(this.writes, values...)
}

func (this *FakeShell) Read(values ...interface{}) {
	for _, value := range values {
		switch query := value.(type) {
		case *queries.EventStream:
			query.Result.Stream = load(this.reads)
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

func (this *FakeShell) PrepareReadResults(results ...interface{}) {
	this.reads = append(this.reads, results)
}

func (this *FakeShell) AssertOutput(expected ...interface{}) {
	this.So(this.writes, should.Resemble, expected)
}
