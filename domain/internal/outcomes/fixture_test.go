package outcomes

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestFixture(t *testing.T) {
    gunit.Run(new(Fixture), t)
}

type Fixture struct {
    *gunit.Fixture
}

func (this *Fixture) Setup() {
}

func (this *Fixture) Test() {
}
