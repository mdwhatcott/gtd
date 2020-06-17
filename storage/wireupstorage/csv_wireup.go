package wireupstorage

import (
	"io"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/storage"
	"github.com/mdwhatcott/gtd/storage/csv"
	"github.com/mdwhatcott/gtd/storage/eventstore"
)

func BuildCachedCSVEventStore(path string) (joyride.StorageReader, joyride.StorageWriter) {
	cache := eventstore.NewCache(
		BuildCSVEventStoreReader(path),
		BuildCSVEventStoreWriter(path),
	)
	return cache, cache
}

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
