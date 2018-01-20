package main

import "github.com/mdwhatcott/gtd/external"

func tasksCLI(input []string) {
	external.Flags(usageFlagsTasks).Parse(input)
	regenerateTasks()
}
