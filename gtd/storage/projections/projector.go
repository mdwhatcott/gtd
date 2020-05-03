package projections

import (
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Projector struct {
	projections map[string]storage.Projection
	eventReader joyride.StorageReader
	eventWriter joyride.StorageWriter
	writer      storage.WriterFunc
}

func NewProjector(
	eventReader joyride.StorageReader,
	eventWriter joyride.StorageWriter,
	writer storage.WriterFunc,
) *Projector {
	return &Projector{
		projections: make(map[string]storage.Projection),
		eventReader: eventReader,
		eventWriter: eventWriter,
		writer:      writer,
	}
}

func (this *Projector) Read(_queries ...interface{}) {
	panic("implement me")
}

func (this *Projector) Write(_events ...interface{}) {
	panic("implement me")
}
