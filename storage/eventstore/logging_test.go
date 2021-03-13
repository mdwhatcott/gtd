package eventstore

import (
	"io"
	"log"

	"github.com/mdwhatcott/gtd/v3/core"
)

func CaptureLogging(out io.Writer) core.Logger {
	return log.New(out, "[CAPTURED] ", 0)
}
