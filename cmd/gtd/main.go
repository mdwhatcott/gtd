package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core/wireupcore"
	"github.com/mdwhatcott/gtd/v3/storage"
	"github.com/mdwhatcott/gtd/v3/storage/wireupstorage"
	"github.com/mdwhatcott/gtd/v3/ui/tempfile"
	"github.com/mdwhatcott/gtd/v3/util/date"
	"github.com/mdwhatcott/gtd/v3/util/exec"
)

var Version = "dev"

func main() {
	log.SetFlags(0)
	log.SetPrefix("  ")
	flag.Usage = PrintUsage
	flag.Parse()

	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable is required for resolution of event store file.")
	}

	PrintBanner(Version)

	APP := BuildApplication(GTDPath)

	ARGS := flag.Args()
	if len(ARGS) == 0 {
		GTDREPL(APP, Version)
	} else {
		GTDOnce(APP, ARGS)
	}

	PushChanges(Version, GTDPath)
}

func PrintBanner(version string) {
	fmt.Println("Welcome to Getting Things Done!")
	fmt.Println()
	fmt.Println("Version:", version)
	fmt.Println()
	fmt.Println("Enter 'help' for instructions.")
}

func BuildApplication(_storageDirectory string) *Application {
	reader, writer := wireupstorage.BuildCachedCSVEventStore(_storageDirectory)
	REQUIREMENTS := wireupcore.Requirements{
		IDFunc: wireupstorage.GenerateID,
		Reader: reader,
		Writer: writer,
	}

	return &Application{
		handler: wireupcore.BuildOutcomesHandler(REQUIREMENTS),
		editor:  tempfile.NewEditor(),
		reader:  REQUIREMENTS.Reader,
	}
}

func PushChanges(_version string, _vcsRoot string) {
	if _version == "dev" {
		return
	}
	STATUS := exec.MustDo(_vcsRoot, "", "git", "status", "--porcelain")
	if !strings.Contains(STATUS, "events.csv") {
		return
	}
	_ = exec.MustDo(_vcsRoot, "Staging newly generated events..." /*****/, "git", "add", storage.EventsDatabaseFilename)
	_ = exec.MustDo(_vcsRoot, "Committing newly generated events..." /**/, "git", "commit", "-m", date.Today())
	_ = exec.MustDo(_vcsRoot, "Pushing newly generated events..." /*****/, "git", "push")

	log.Println("Finished.")
}
