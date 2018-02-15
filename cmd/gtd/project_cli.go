package main

import "github.com/mdwhatcott/gtd/external"

func projectCLI(input []string) {
	flag := external.Flags(usageFlagsProject)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "":
		createManyProjectsCLI(remaining)
	case "create":
		createProjectCLI(remaining)
	case "update":
		updateProjectStatusCLI(remaining)
	case "list":
		listProjectsCLI(remaining)
	default:
		exit(flag)
	}
}
