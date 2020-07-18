package projections

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core"
)

type ProjectorFixture struct {
	*gunit.Fixture
	projector core.Projector
}

func InitializeProjectorFixture(_inner *gunit.Fixture, _projector core.Projector) *ProjectorFixture {
	return &ProjectorFixture{Fixture: _inner, projector: _projector}
}
func (this *ProjectorFixture) apply(_events ...interface{}) {
	STREAM := make(chan interface{}, len(_events))
	for _, EVENT := range _events {
		STREAM <- EVENT
	}
	close(STREAM)
	this.projector.Apply(STREAM)
}
func (this *ProjectorFixture) assert(_expected interface{}) {
	this.So(this.projector.Projection(), should.Resemble, _expected)
}
