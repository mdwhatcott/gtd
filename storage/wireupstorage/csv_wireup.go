package wireupstorage

import (
	"io"
	"path/filepath"

	"github.com/smartystreets/joyride/v2"

	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/storage/csv"
	"github.com/mdwhatcott/gtd/v3/storage/eventstore"
)

func BuildCachedCSVEventStore(folder string) (joyride.StorageReader, joyride.StorageWriter) {
	PATH := filepath.Join(folder, storage.EventsDatabaseFilename)
	CACHE := eventstore.NewCache(
		BuildCSVEventStoreReader(PATH),
		BuildCSVEventStoreWriter(PATH),
	)
	return CACHE, CACHE
}

func BuildCSVEventStoreReader(path string) joyride.StorageReader {
	return eventstore.NewReader(reading(path), csvDecoding)
}
func BuildCSVEventStoreWriter(path string) joyride.StorageWriter {
	return eventstore.NewWriter(csvEncoding, writing(path))
}

func csvDecoding(reader io.Reader) storage.Decoder {
	return csv.NewDecoder(reader, csv.DecoderRegistry())
}
func csvEncoding(writer io.Writer) storage.Encoder {
	return csv.NewEncoder(writer, csv.EncoderRegistry())
}
