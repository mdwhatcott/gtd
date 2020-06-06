package projections

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/core"
)

type ProjectorFixture struct {
	*gunit.Fixture
	projector core.Projector
}

func InitializeProjectorFixture(_inner *gunit.Fixture, _projector core.Projector) *ProjectorFixture {
	return &ProjectorFixture{Fixture: _inner, projector: _projector}
}
func (this *ProjectorFixture) apply(_events ...interface{}) {
	this.projector.Apply(_events...)
}
func (this *ProjectorFixture) assert(_expected interface{}) {
	this.So(this.projector.Projection(), should.Resemble, _expected)
}
