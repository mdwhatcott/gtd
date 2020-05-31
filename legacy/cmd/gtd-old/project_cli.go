package main

import "github.com/mdwhatcott/gtd/legacy/external"

func projectCLI(input []string) {
	flag := external.Flags(usageFlagsProject)
	_ = flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "":
		createProjects()
		regenerateTasks()
	case "update":
		updateProjectStatusCLI(remaining)
		regenerateTasks()
	case "list":
		listProjects()
	case "review":
		reviewProjects()
		regenerateTasks()
	default:
		exit(flag)
	}
}
