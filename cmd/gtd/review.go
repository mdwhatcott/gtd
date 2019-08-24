package main

import (
	"fmt"
	"time"

	"github.com/mdwhatcott/tomato"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func review() {
	_ = external.Flags(usageFlagsReview).Parse(nil)

	getClear()
	getCurrent()
	getCreative()
}

func getClear() {
	fmt.Println("Part 1: Get Clear")
	commit()
	mindSweep()
	inboxZero()
	calendar()
}

func getCurrent() {
	fmt.Println("Part 2: Get Current")
	importTickler()
	generateTasks()
	taskReview()
	generateTasks()
	projectReview()
	generateTasks()
	commit()
}
func getCreative() {
	fmt.Println("Part 3: Get Creative! (review someday/maybe, upper horizons of focus, etc...)")
}

func commit() {
	fmt.Println("Commit current state. <ENTER> to continue...")
	external.Commit(gtd.FolderRoot)
}

func mindSweep() {
	const durationPrompt = "What has your attention? Enter duration of mind sweep (or <ENTER> to continue):"
	answer := external.Prompt(durationPrompt)
	for {
		if len(answer) == 0 {
			break
		} else if duration, err := time.ParseDuration(answer); err != nil {
			answer = external.Prompt(durationPrompt)
			continue
		} else {
			tomato.SetTimer(duration).Start()
			break
		}
	}
}

const enterToContinue = " <ENTER> to continue..."

// TODO: promote this to a CLI action
func inboxZero() {
	fmt.Println(`Get "In" to Zero!`)
	external.Prompt("1. Clear desk, moving 'stuff' to the in-tray." + enterToContinue)
	external.Prompt("2. Unload stuff from messenger bag to the in-tray." + enterToContinue)
	external.Prompt("3. Process in-tray." + enterToContinue)
	createProjects()
	external.Prompt("4. Process ~/Desktop." + enterToContinue)
	createProjects()
	external.Prompt("5. Process Google Keep." + enterToContinue)
	createProjects()
	external.Prompt("6. Process e-mail." + enterToContinue)
	createProjects()
	external.Prompt("All in-boxes should be empty at this point. <ENTER> to continue...")
}

func calendar() {
	external.Prompt("Review previous and upcoming 2 weeks in calendar." + enterToContinue)
	createProjects()
}

// TODO: should this happen any time we 'generateTasks'?
func importTickler() {
	fmt.Println("Scanning tickler for projects now due...")
	scanTickler()
}

func generateTasks() {
	fmt.Println("Regenerating task lists...")
	regenerateTasks()
}

func taskReview() {
	fmt.Println("Review task lists in sequence...")
	reviewTasks()
}

func projectReview() {
	fmt.Println("Review project status in sequence.")
	reviewProjects()
}
