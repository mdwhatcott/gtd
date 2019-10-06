package main

import (
	"log"
	"os"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func main() {
	flag := external.Flags(usageFlag)
	_ = flag.Parse(os.Args[1:])

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "inbox":
		inboxZero()
	case "sweep":
		mindSweep()
	case "review":
		review()
	case "project":
		projectCLI(remaining)
	case "tasks":
		tasksCLI(remaining)
	case "tickler":
		scanTickler()
	case "notes":
		external.OpenTextEditor(gtd.NotesRoot)
	case "":
		external.OpenTextEditor(gtd.FolderRoot)
	default:
		log.Println("Unknown args:", flag.Args())
		flag.Usage()
		os.Exit(1)
	}
}
