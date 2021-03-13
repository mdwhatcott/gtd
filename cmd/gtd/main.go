package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core"
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
	PullLatest(Version, GTDPath)

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
	fmt.Println()
}

func BuildApplication(storageDirectory string) *Application {
	reader, writer := wireupstorage.BuildCachedCSVEventStore(log.Default(), storageDirectory)
	REQUIREMENTS := wireupcore.Requirements{
		Log:    log.Default(),
		Clock:  core.Now,
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

func PullLatest(version string, vcsRoot string) {
	if version == "dev" {
		return
	}
	STATUS := exec.MustDo(vcsRoot, "", "git", "status", "--porcelain")
	if strings.Contains(STATUS, "events.csv") {
		log.Fatal("Events database is in 'dirty' state. Please commit and push changes manually and restart.")
		return
	}
	_ = exec.MustDo(vcsRoot, "Fetching changes...", "git", "fetch")
	_ = exec.MustDo(vcsRoot, "Rebasing changes...", "git", "rebase", "origin/main")

	log.Println("Ready.")
}
func PushChanges(version string, vcsRoot string) {
	if version == "dev" {
		return
	}
	STATUS := exec.MustDo(vcsRoot, "", "git", "status", "--porcelain")
	if !strings.Contains(STATUS, "events.csv") {
		return
	}
	log.Println("Total events:", exec.MustDo(vcsRoot, "", "wc", "-l", "events.csv"))
	_ = exec.MustDo(vcsRoot, "Staging newly generated events..." /*****/, "git", "add", storage.EventsDatabaseFilename)
	_ = exec.MustDo(vcsRoot, "Committing newly generated events..." /**/, "git", "commit", "-m", date.Today())
	_ = exec.MustDo(vcsRoot, "Pushing newly generated events..." /*****/, "git", "push")

	log.Println("Finished.")
}
