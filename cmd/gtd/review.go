package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func reviewCLI(inputs []string) {
	external.Flags(usageFlagsReview).Parse(inputs)

	fmt.Println("Step 0/4: Commit all changes (this will happen serveral more times along the way).")
	external.Commit(gtd.FolderRoot)

	fmt.Println("Step 1/4: Process inbox, creating new projects as needed.")
	createManyProjects()
	external.Commit(gtd.FolderRoot)

	fmt.Println("Step 2/4: Synchronize task completion status from context lists to project documents.")
	syncTasks(gtd.LoadProjects())
	external.Commit(gtd.FolderRoot)

	fmt.Println("Step 3/4: Review each project's content, tasks, and status.")
	listProjects(gtd.LoadProjects(), true)
	external.Commit(gtd.FolderRoot)

	fmt.Println("Step 4/4: Regenerate contextualized task lists.")
	sweepTasks(gtd.LoadProjects())
	external.Commit(gtd.FolderRoot)

	fmt.Println("Review procedure complete.")
}
