package eventstore

import (
	"io"

	"github.com/smartystreets/logging"
)

func CaptureLogging(out io.Writer) *logging.Logger {
	log := logging.Capture()
	log.SetFlags(0)
	log.SetOutput(out)
	log.SetPrefix("[CAPTURED] ")
	return log
}
