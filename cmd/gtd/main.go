package main

import "os"

func main() {
	flag := flags(usageFlag)
	flag.Parse(os.Args[1:])

	first, remaining := firstAndRemaining(flag.Args())

	switch first {
	case "review":
		weeklyReview(remaining)
	case "project":
		parseProjectCLI(remaining)
	case "tasks":
		parseTaskCLI(remaining)
	default:
		exit(flag)
	}
}
