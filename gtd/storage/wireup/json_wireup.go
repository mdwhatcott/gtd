package wireup

import (
	"io"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
	"github.com/mdwhatcott/gtd/gtd/storage/json"
)

func BuildJSONEventStoreReader(path string) joyride.StorageReader {
	return eventstore.NewReader(reading(path), jsonDecoding)
}
func BuildJSONEventStoreWriter(path string) joyride.StorageWriter {
	return eventstore.NewWriter(jsonEncoding, writing(path))
}

func jsonDecoding(_reader io.Reader) storage.Decoder {
	return json.NewDecoder(_reader, json.Registry())
}
func jsonEncoding(_writer io.Writer) storage.Encoder { return json.NewEncoder(_writer) }
