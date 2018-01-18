package main

func parseProjectCLI(input []string) {
	flag := flags(usageFlagsProject)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "list":
		listProjects(remaining)
	case "create":
		createProject(remaining)
	case "create-many":
		createManyProjects(remaining)
	case "update":
		updateProject(remaining)
	default:
		exit(flag)
	}
}
