package wireup

import (
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/domain"
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
