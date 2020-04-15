package domain

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture

	shell *FakeShell
}

func (this *Fixture) Setup() {
	this.shell = NewFakeShell(this.Fixture)
}

func (this *Fixture) Test() {
}
