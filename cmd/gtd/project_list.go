package main

import (
	"fmt"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func listProjectsCLI(inputs []string) {
	flag := external.Flags(usageFlagsListProjects)
	review := flag.Bool("review", false, "When set, review each project via a REPL and text editor sessions.")
	flag.Parse(inputs)
	projects := gtd.LoadProjects()
	listProjects(projects, *review)
}

func listProjects(projects []*gtd.Project, review bool) {
	for _, project := range projects {
		fmt.Println(project.Name())
		if review {
			external.OpenTextEditorAndWait(project.Path())
			interactiveUpdateProjectStatus(project)
		}
	}
}

func interactiveUpdateProjectStatus(project *gtd.Project) {
	for {
		fmt.Print("Enter new project status or <blank> to skip: (complete, someday, maybe, reject) ")
		if line := external.ReadLine(); line == "" {
			break
		} else if updateProjectStatus(project, line) {
			break
		} else {
			fmt.Println("Invalid status:", line)
		}
	}
}
