package domain

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/joyride/v2"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture

	shell   *FakeShell
	handler *Handler
}

func (this *Fixture) Setup() {
	this.shell = NewFakeShell(this.Fixture)
}

func (this *Fixture) handle(commands ...interface{}) {
	this.handler = NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(this.shell),
			joyride.WithStorageWriter(this.shell),
			joyride.WithMessageDispatcher(this.shell),
		),
	)
	this.handler.Handle(commands...)
}

func (this *Fixture) TestHandlerPanicsOnUnrecognizedMessageTypes() {
	this.So(func() { this.handle(42) }, should.PanicWith, joyride.ErrUnknownType)
	this.So(func() { this.handle(true) }, should.PanicWith, joyride.ErrUnknownType)
}
