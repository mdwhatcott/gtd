package main

import (
	"os"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func main() {
	flag := external.Flags(usageFlag)
	flag.Parse(os.Args[1:])

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "review":
		reviewCLI(remaining)
	case "project":
		projectCLI(remaining)
	case "tasks":
		tasksCLI(remaining)
	default:
		external.OpenTextEditor(gtd.FolderRoot)
	}
}
