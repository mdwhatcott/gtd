package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func reviewCLI(inputs []string) {
	external.Flags(usageFlagsReview).Parse(inputs)
	external.Commit(gtd.FolderRoot)
	regenerateTasks()
	createProjects()
	reviewProjects()
	regenerateTasks()
	external.Commit(gtd.FolderRoot)
	fmt.Println("Finished.")
}
