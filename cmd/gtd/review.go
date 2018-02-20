package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func reviewCLI(inputs []string) {
	external.Flags(usageFlagsReview).Parse(inputs)

	fmt.Println("Step 0/4: Commit and push all changes now (sourcetree window will open shortly).")
	external.Commit(gtd.FolderRoot)

	fmt.Println("Step 1/4: Process inbox, creating new projects as needed.")
	createManyProjects()

	fmt.Println("Step 2/4: Regenerate all task lists.")
	syncTasks()

	fmt.Println("Step 3/4: Review each project's content, tasks, and status.")
	listProjects(gtd.LoadProjects(), true)

	fmt.Println("Step 4/4: Regenerate all task lists.")
	regenerateTasks()

	fmt.Println("Review procedure complete. Commit and push all changes now.")
	external.Commit(gtd.FolderRoot)

}
