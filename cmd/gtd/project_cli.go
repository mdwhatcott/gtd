package main

import "github.com/mdwhatcott/gtd/external"

func projectCLI(input []string) {
	flag := external.Flags(usageFlagsProject)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "":
		createProjects()
	case "update":
		updateProjectStatusCLI(remaining)
	case "list":
		listProjects()
	case "review":
		reviewProjects()
	default:
		exit(flag)
	}
}
