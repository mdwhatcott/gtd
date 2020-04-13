package outcomes

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/joyride"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture

	shell  *FakeShell
	runner joyride.Runner
}

func (this *Fixture) Setup() {
	this.shell = NewFakeShell(this.Fixture)
	this.runner = joyride.NewRunner(this.shell, this.shell, this.shell)
}

func (this *Fixture) Test() {
}

///////////////////////////////////////////////////////////

type FakeShell struct {
	*gunit.Fixture
	writes   []interface{}
	messages []interface{}
	reads    map[string][]interface{} // not really sure about this
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
		switch value.(type) {
		// TODO: fill out results based on type of value
		}
	}
}

func (this *FakeShell) PrepareReadResults(results ...interface{}) {
	// TODO: populate reads...
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
