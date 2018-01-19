package main

import "github.com/mdwhatcott/gtd/external"

func taskCLI(input []string) {
	flag := external.Flags(usageFlagsTasks)
	flag.Parse(input)

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "sync":
		syncTasksCLI(remaining)
	case "sweep":
		sweepTasksCLI(remaining)
	default:
		exit(flag)
	}
}
