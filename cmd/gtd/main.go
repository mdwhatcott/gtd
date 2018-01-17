package main

import "flag"

func main() {
	flag.Parse()

	switch args := append(flag.Args(), ""); args[0] {
	case "review":
		weeklyReview(args[1:])
	case "project":
		parseProjectCLI(args[1:])
	case "task":
		parseTaskCLI(args[1:])
	default:
		flag.Usage()
	}
}
