package wireupcore

import (
	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/domain"
)

type Requirements struct {
	IDFunc core.IDFunc
	Reader joyride.StorageReader
	Writer joyride.StorageWriter
}

func BuildOutcomesHandler(_requirements Requirements) core.Handler {
	return domain.NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(_requirements.Reader),
			joyride.WithStorageWriter(_requirements.Writer),
		),
		_requirements.IDFunc,
	)
}
