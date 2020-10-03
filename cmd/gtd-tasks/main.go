package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core/projections"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/storage/wireupstorage"
	"github.com/mdwhatcott/gtd/v3/ui/ux"
)

func main() {
	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable is required for resolution of event store file.")
	}
	PATH := filepath.Join(GTDPath, storage.EventsDatabaseFilename)
	READER := wireupstorage.BuildCSVEventStoreReader(PATH)
	PROJECTION := replayIncompleteActions(READER)
	MARKDOWN := ux.FormatIncompleteActions(PROJECTION.Contexts...)
	fmt.Println(MARKDOWN)
}

func replayIncompleteActions(reader joyride.StorageReader) projections.IncompleteActionsByContext {
	STREAM := &storage.EventStream{}
	reader.Read(context.Background(), STREAM)
	PROJECTOR := projections.NewIncompleteActionsByContextProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
	return PROJECTOR.IncompleteActionsByContextProjection()
}

