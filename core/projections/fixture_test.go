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

func InitializeProjectorFixture(inner *gunit.Fixture, projector core.Projector) *ProjectorFixture {
	return &ProjectorFixture{Fixture: inner, projector: projector}
}
func (this *ProjectorFixture) apply(events ...interface{}) {
	STREAM := make(chan interface{}, len(events))
	for _, EVENT := range events {
		STREAM <- EVENT
	}
	close(STREAM)
	this.projector.Apply(STREAM)
}
func (this *ProjectorFixture) assert(expected interface{}) {
	this.So(this.projector.Projection(), should.Resemble, expected)
}
