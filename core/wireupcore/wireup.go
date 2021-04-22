package wireupcore

import (
	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/domain"
)

type Requirements struct {
	Log    core.Logger
	Clock  core.Clock
	IDFunc core.IDFunc
	Reader joyride.StorageReader
	Writer joyride.StorageWriter
}

func BuildOutcomesHandler(_requirements Requirements) core.Handler {
	return domain.NewHandler(
		_requirements.Log,
		_requirements.Clock,
		_requirements.IDFunc,
		joyride.NewRunner(
			joyride.WithStorageReader(_requirements.Reader),
			joyride.WithStorageWriter(_requirements.Writer),
		),
	)
}
