package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
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
