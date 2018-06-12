package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func review() {
	external.Flags(usageFlagsReview).Parse(nil)
	external.Commit(gtd.FolderRoot)
	fmt.Println("Scanning tickler...")
	scanTickler()
	fmt.Println("Regenerating task lists...")
	regenerateTasks()
	fmt.Println("Creating projects...")
	createProjects()
	fmt.Println("Reviewing projects...")
	reviewProjects()
	fmt.Println("Regenerating task lists...")
	regenerateTasks()
	external.Commit(gtd.FolderRoot)
	fmt.Println("Finished.")
}
