package main

import "github.com/mdwhatcott/gtd/external"

func tasks(input []string) {
	external.Flags(usageFlagsTasks).Parse(input)
	syncTasks()
	sweepTasks()
}
