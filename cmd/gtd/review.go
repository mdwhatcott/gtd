package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func review() {
	external.Flags(usageFlagsReview).Parse(nil)
	external.Commit(gtd.FolderRoot)
	regenerateTasks()
	createProjects()
	reviewProjects()
	regenerateTasks()
	external.Commit(gtd.FolderRoot)
	fmt.Println("Finished.")
}
