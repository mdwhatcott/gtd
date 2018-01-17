package main

func parseProjectCLI(input []string) {
	projectFlags.Parse(input)

	args := projectFlags.Args()
	switch args[0] {
	case "list":
		listProjects(input[1:])
	case "create":
		createProject(input[1:])
	case "create-many":
		createManyProjects(input[1:])
	case "update":
		updateProject(input[1:])
	default:
		projectFlags.Usage()
	}
}
