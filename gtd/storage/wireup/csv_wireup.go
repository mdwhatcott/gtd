package wireup

import (
	"io"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/gtd/storage"
	"github.com/mdwhatcott/gtd/gtd/storage/csv"
	"github.com/mdwhatcott/gtd/gtd/storage/eventstore"
)

func BuildCSVEventStoreReader(path string) joyride.StorageReader {
	return eventstore.NewReader(reading(path), csvDecoding)
}
func BuildCSVEventStoreWriter(path string) joyride.StorageWriter {
	return eventstore.NewWriter(csvEncoding, writing(path))
}

func csvDecoding(_reader io.Reader) storage.Decoder {
	return csv.NewDecoder(_reader, csv.DecoderRegistry())
}
func csvEncoding(_writer io.Writer) storage.Encoder {
	return csv.NewEncoder(_writer, csv.EncoderRegistry())
}
