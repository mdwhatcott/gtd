package main

import (
	"fmt"
	"time"

	"github.com/mdwhatcott/gtd/legacy/external"
	"github.com/mdwhatcott/gtd/legacy/gtd"
)

func listProjects() {
	for _, project := range gtd.LoadProjects() {
		fmt.Println(project.Name())
	}
}

func reviewProjects() {
	time.Sleep(time.Second)
	for _, project := range gtd.LoadProjects() {
		fmt.Println("Reviewing:", project.Name())
		external.OpenTextEditorAndWait(project.Path())
		interactiveUpdateProjectStatus(project)
	}
}

func interactiveUpdateProjectStatus(project *gtd.Project) {
	for {
		fmt.Print("- Enter new project status or <blank> to skip: (complete, someday, maybe, reject) ")
		if line := external.ReadLine(); line == "" {
			break
		} else if updateProjectStatus(project, line) {
			break
		} else {
			fmt.Println("Invalid status:", line)
		}
	}
}
