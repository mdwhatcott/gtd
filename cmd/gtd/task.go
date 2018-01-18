package main

func parseTaskCLI(input []string) {
	flag := flags(usageFlagsTasks)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "sync":
		sweepTasks(remaining)
	case "sweep":
		syncTasks(remaining)
	default:
		exit(flag)
	}
}
