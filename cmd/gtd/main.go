package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdwhatcott/gtd/v3/core/wireupcore"
	"github.com/mdwhatcott/gtd/v3/storage/wireupstorage"
	"github.com/mdwhatcott/gtd/v3/ui/tempfile"
)

var Version = "dev"

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Usage = PrintUsage
	flag.Parse()

	GTDPath, OK := os.LookupEnv("GTDPATH")
	if !OK {
		log.Fatal("The 'GTDPATH' environment variable required for resolution of event store file.")
	}

	PrintBanner(Version)

	APP := BuildApplication(GTDPath)

	ARGS := flag.Args()
	if len(ARGS) == 0 {
		DoREPL(APP, Version)
	} else {
		DoOnce(APP, ARGS)
	}

	PushChanges(GTDPath)
}

func BuildApplication(storageDirectory string) *Application {
	PATH := filepath.Join(storageDirectory, "events.csv")
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

		storageDirectory: storageDirectory,
	}
}

func PushChanges(storageDirectory string) {
	STATUS := exec.Command("git", "status", "--porcelain")
	STATUS.Dir = storageDirectory
	OUT, ERR := STATUS.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	if !strings.Contains(strings.TrimSpace(string(OUT)), "events.csv") {
		return
	}

	log.Println("Staging newly generated events...")
	ADD := exec.Command("git", "add", "events.csv")
	ADD.Dir = storageDirectory
	OUT, ERR = ADD.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	log.Println("Committing newly generated events...")
	TODAY := time.Now().Format("2006-01-02")
	COMMIT := exec.Command("git", "commit", "-m", TODAY)
	COMMIT.Dir = storageDirectory
	OUT, ERR = COMMIT.CombinedOutput()
	if ERR != nil {
		log.Println(OUT)
		log.Fatal(ERR)
	}

	log.Println("Pushing newly generated events...")
	PUSH := exec.Command("git", "push", "origin", "main")
	PUSH.Dir = storageDirectory
	OUT, ERR = PUSH.CombinedOutput()
	if ERR != nil {
		log.Println(PUSH)
		log.Fatal(ERR)
	}

	log.Println("Finished.")
}
