package wireup

import (
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/core/events"
	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
)

func BuildEventStoreReader(path string) joyride.StorageReader {
	return eventstore.NewReader(reading(path), decoding)
}
func BuildEventStoreWriter(path string) joyride.StorageWriter {
	return eventstore.NewWriter(encoding, writing(path))
}

func decoding(_reader io.Reader) storage.Decoder { return json.NewDecoder(_reader, events.Registry()) }
func encoding(_writer io.Writer) storage.Encoder { return json.NewEncoder(_writer) }

func reading(path string) func() io.ReadCloser  { return func() io.ReadCloser { return open(path) } }
func writing(path string) func() io.WriteCloser { return func() io.WriteCloser { return open(path) } }

func open(_path string) *os.File {
	FILE, ERR := os.OpenFile(_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if ERR != nil {
		panic(ERR)
	}
	return FILE
}

func GenerateID() string {
	return uuid.New().String()
}
