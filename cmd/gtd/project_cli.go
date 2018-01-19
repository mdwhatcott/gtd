package main

import "github.com/mdwhatcott/gtd/external"

func projectCLI(input []string) {
	flag := external.Flags(usageFlagsProject)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "list":
		listProjectsCLI(remaining)
	case "create":
		createProjectCLI(remaining)
	case "create-many":
		createManyProjectsCLI(remaining)
	case "update":
		updateProjectStatusCLI(remaining)
	default:
		exit(flag)
	}
}
