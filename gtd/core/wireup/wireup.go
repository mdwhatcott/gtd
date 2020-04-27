package wireup

import (
	"github.com/google/uuid"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/internal/outcomes"
)

func BuildHandler(reader joyride.StorageReader, writer joyride.StorageWriter) core.Handler {
	return outcomes.NewHandler(
		joyride.NewRunner(
			joyride.WithStorageReader(reader),
			joyride.WithStorageWriter(writer),
		),
		outcomes.NewTask(idFunc),
	)
}

func idFunc() string {
	return uuid.New().String()
}
