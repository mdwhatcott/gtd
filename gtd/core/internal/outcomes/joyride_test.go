package outcomes

import (
	"github.com/smartystreets/assertions/should"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type FakeShell struct {
	*Fixture

	writes []interface{}
	reads  map[string][]interface{}
}

func NewFakeShell(_fixture *Fixture) *FakeShell {
	return &FakeShell{
		Fixture: _fixture,
		reads:   make(map[string][]interface{}),
	}
}

func (this *FakeShell) Write(_values ...interface{}) {
	this.writes = append(this.writes, _values...)
}

func (this *FakeShell) Read(_values ...interface{}) {
	this.log.Println("Reading:", _values)
	for _, VALUE := range _values {
		this.log.Println("Reading value:", VALUE)
		switch QUERY := VALUE.(type) {
		case *storage.OutcomeEventStream:
			this.log.Println("Reading outcome event stream...", QUERY.OutcomeID)
			QUERY.Result.Stream = make(chan interface{})
			go load(QUERY.Result.Stream, this.reads[QUERY.OutcomeID])
		}
	}
}

func load(_stream chan interface{}, _events []interface{}) {
	defer close(_stream)
	for _, EVENT := range _events {
		_stream <- EVENT
	}
}

func (this *FakeShell) PrepareReadResults(_id string, _results ...interface{}) {
	this.reads[_id] = append(this.reads[_id], _results...)
	this.log.Println("Read:", _id, _results)
}

func (this *FakeShell) AssertNoOutput() {
	this.AssertOutput()
}

func (this *FakeShell) AssertOutput(_expected ...interface{}) {
	this.So(this.writes, should.Resemble, _expected)
}
