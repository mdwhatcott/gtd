package main

import (
	"fmt"

	"github.com/mdwhatcott/tomato"

	"github.com/mdwhatcott/gtd/legacy/external"
	"github.com/mdwhatcott/gtd/legacy/gtd"
)

func review() {
	_ = external.Flags(usageFlagsReview).Parse(nil)

	getClear()
	getCurrent()
	getCreative()
}

func getClear() {
	fmt.Println("Part 1: Get Clear")
	mindSweep()
	inboxZero()
	commit()
}
func mindSweep() {
	duration := external.PromptDuration(durationPrompt)
	if duration > 0 {
		tomato.SetTimer(duration).Start()
	}
}
func inboxZero() {
	fmt.Println(`Get "In" to Zero!`)
	clearDesk()
	unloadMessengerBag()
	processInTray()
	processDesktop()
	processEmail()
	external.Prompt("All in-boxes should be empty at this point. <ENTER> to continue...")
}
func clearDesk() string {
	return external.Prompt("Clear the desk, moving 'stuff' to the in-tray." + enterToContinue)
}
func unloadMessengerBag() string {
	return external.Prompt("Unload stuff from messenger bag to the in-tray." + enterToContinue)
}
func processInTray() {
	external.Prompt("Process the in-tray." + enterToContinue)
	createProjects()
}
func processDesktop() {
	external.Prompt("Process ~/Desktop." + enterToContinue)
	createProjects()
}
func processEmail() {
	external.Prompt("Process e-mail." + enterToContinue)
	createProjects()
}

func getCurrent() {
	fmt.Println("Part 2: Get Current")
	reviewTextMessages()
	reviewEmailTags()
	reviewCalendar()
	importTickler()
	taskReview()
	projectReview()
	reviewProjectList()
	commit()
}
func reviewTextMessages() {
	external.Prompt("Review recent text messages." + enterToContinue)
	createProjects()
}
func reviewEmailTags() {
	external.Prompt("Review email tags (needs reply, read later, project support, waiting for response)." + enterToContinue)
}
func reviewCalendar() {
	external.Prompt("Review previous and upcoming 2 weeks in calendar." + enterToContinue)
	createProjects()
}
func importTickler() {
	fmt.Println("Scanning tickler for projects now due...")
	scanTickler()
	generateTasks()
}
func taskReview() {
	fmt.Println("Review task lists in sequence...")
	reviewTasks()
	generateTasks()
}
func projectReview() {
	fmt.Println("Review project status in sequence.")
	reviewProjects()
	generateTasks()
}
func reviewProjectList() {
	fmt.Println("Here is the projects list for one final review:")
	listProjects()
}
func generateTasks() {
	fmt.Println("Regenerating task lists...")
	regenerateTasks()
}

func getCreative() {
	fmt.Println("Part 3: Get Creative! ('clean a drawer', review someday/maybe, review upper horizons of focus, etc...)")
}

func commit() {
	fmt.Println("Commit current state. <ENTER> to continue...")
	external.Commit(gtd.FolderRoot)
}

const enterToContinue = " <ENTER> to continue..."
const durationPrompt = "What has your attention? Enter duration of mind sweep (or <ENTER> to continue):"
