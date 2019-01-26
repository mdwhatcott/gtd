package main

import (
	"fmt"
	"time"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
	"github.com/mdwhatcott/tomato"
)

func review() {
	_ = external.Flags(usageFlagsReview).Parse(nil)

	getClear()
	getCurrent()
	getCreative()
}

func getClear() {
	fmt.Println("Part 1: 'Get Clear'")
	commit()
	mindSweep()
	inboxZero()
}
func getCurrent() {
	fmt.Println("Part 2: 'Get Current'")
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
	answer := external.Prompt("What has your attention? Enter duration of mind sweep:")
	for {
		duration, err := time.ParseDuration(answer)
		if err != nil {
			answer = external.Prompt("What has your attention? Enter duration of mind sweep:")
			continue
		} else {
			tomato.SetTimer(duration).Start()
			break
		}
	}
}

func inboxZero() {
	fmt.Println(`Get "In" to Zero:
				
1. Gather all physical inputs into in-basket.
2. Process all items in Google Keep.
3. Process all emails.
4. Review previous and upcoming 2 weeks in calendar.
			
(Browser tabs and project creation prompt will appear shortly...)`)

	time.Sleep(time.Second * 5)
	external.Navigate("https://keep.google.com")
	external.Navigate("https://mail.google.com")
	external.Navigate("https://calendar.google.com")
	createProjects()
	external.Prompt("All inboxes should be empty at this point. <ENTER> to continue...")
}

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
