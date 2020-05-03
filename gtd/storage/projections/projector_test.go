package projections

import (
	"testing"

	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/util/fake"
)

func TestProjectorFixture(t *testing.T) {
	gunit.Run(new(ProjectorFixture), t)
}

type ProjectorFixture struct {
	*gunit.Fixture
	eventReader *fake.Reader
	eventWriter *fake.Writer
	projector   *Projector
}

func (this *ProjectorFixture) Setup() {
	this.projector = NewProjector(nil, nil, nil)
}

func (this *ProjectorFixture) Test() {
}
