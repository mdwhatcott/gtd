package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func createManyProjectsCLI(inputs []string) {
	external.Flags(usageFlagsCreateManyProjects).Parse(inputs)
	createManyProjects()
}

func createManyProjects() {
	for {
		fmt.Print("Enter project name (<blank> to quit): ")
		name := external.ReadLine()
		if name == "" {
			break
		}
		createProject(gtd.CreateProjectCommand{Name: name})
	}
}
