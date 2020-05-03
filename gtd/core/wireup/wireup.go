package wireup

import (
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/internal/outcomes"
)

type Requirements struct {
	IDFunc func() string
	Reader joyride.StorageReader
	Writer joyride.StorageWriter
}

func BuildHandler(requirements Requirements) core.Handler {
	return outcomes.NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(requirements.Reader),
			joyride.WithStorageWriter(requirements.Writer),
		),
		outcomes.NewTask(requirements.IDFunc),
	)
}
