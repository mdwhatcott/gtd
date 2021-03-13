package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartystreets/joyride/v3"

	"github.com/mdwhatcott/gtd/v3/core/projections"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/storage/wireupstorage"
)

var Version = "dev"

func main() {
	log.SetFlags(0)
	log.Println("gtd-actions@" + Version)

	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable is required for resolution of event store file.")
	}
	PATH := filepath.Join(GTDPath, storage.EventsDatabaseFilename)
	READER := wireupstorage.BuildCSVEventStoreReader(nil, PATH)
	PROJECTION := replayIncompleteActions(READER)
	MARKDOWN := FormatIncompleteActions(PROJECTION.Contexts...)
	fmt.Println("# Actions - " + time.Now().Format("2006-01-02"))
	fmt.Println(MARKDOWN)
}

func replayIncompleteActions(reader joyride.StorageReader) projections.IncompleteActionsByContext {
	STREAM := &storage.EventStream{}
	reader.Read(context.Background(), STREAM)
	PROJECTOR := projections.NewIncompleteActionsByContextProjector()
	PROJECTOR.Apply(STREAM.Result.Events)
	return PROJECTOR.IncompleteActionsByContextProjection()
}

func FormatIncompleteActions(contexts ...*projections.Context) string {
	BUILDER := new(strings.Builder)

	for _, CONTEXT := range contexts {
		_, _ = fmt.Fprintf(BUILDER, "## @%s:\n\n", strings.Title(CONTEXT.Name))

		for _, ACTION := range CONTEXT.Actions {
			_, _ = fmt.Fprintf(BUILDER,
				"- [ ] %s (%s)\n",
				ACTION.Description,
				ACTION.OutcomeTitle,
			)
		}

		BUILDER.WriteString("\n\n")
	}

	return strings.TrimSpace(BUILDER.String())
}
