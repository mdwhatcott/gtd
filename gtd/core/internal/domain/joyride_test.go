package domain

import (
	"github.com/smartystreets/assertions/should"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type FakeShell struct {
	*Fixture

	writes []interface{}
	reads  map[string][]interface{}
}

func NewFakeShell(fixture *Fixture) *FakeShell {
	return &FakeShell{
		Fixture: fixture,
		reads:   make(map[string][]interface{}),
	}
}
func (this *FakeShell) Write(values ...interface{}) {
	this.writes = append(this.writes, values...)
}
func (this *FakeShell) Read(values ...interface{}) {
	for _, value := range values {
		switch query := value.(type) {
		case *storage.OutcomeEventStream:
			go this.populate(query.OutcomeID, query.Result.Stream)
		}
	}
}
func (this *FakeShell) populate(id string, stream chan interface{}) {
	for _, event := range this.reads[id] {
		//this.Printf("adding event to stream: %#v", event)
		stream <- event
	}
	close(stream)
}
func (this *FakeShell) PrepareReadResults(id string, results ...interface{}) {
	this.reads[id] = append(this.reads[id], results...)
}
func (this *FakeShell) AssertNoOutput() {
	this.AssertOutput()
}
func (this *FakeShell) AssertOutput(expected ...interface{}) {
	this.So(this.writes, should.Resemble, expected)
}
