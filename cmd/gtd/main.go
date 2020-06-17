package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/gtd/core/wireupcore"
	"github.com/mdwhatcott/gtd/storage/wireupstorage"
	"github.com/mdwhatcott/gtd/ui/tempfile"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Usage = Usage
	flag.Parse()

	REPL(BuildApplication(), flag.Args())
}

func BuildApplication() *Application {
	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable required for resolution of event store file.")
	}

	PATH := filepath.Join(GTDPath, "events.csv")
	reader, writer := wireupstorage.BuildCachedCSVEventStore(PATH)
	REQUIREMENTS := wireupcore.Requirements{
		IDFunc: wireupstorage.GenerateID,
		Reader: reader,
		Writer: writer,
	}

	return &Application{
		handler: wireupcore.BuildOutcomesHandler(REQUIREMENTS),
		editor:  tempfile.NewEditor(),
		reader:  REQUIREMENTS.Reader,

		storageDirectory: GTDPath,
	}
}
