package main

import (
	"fmt"
	"log"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func main() {
	projects := gtd.LoadProjects()
	for context, tasks := range sortTasksByContext(projects) {
		fmt.Println("## " + context)
		fmt.Println()
		for _, task := range tasks {
			fmt.Println(task.PrintableString())
		}
		fmt.Println()
	}

}

func sortTasksByContext(projects []*gtd.Project) map[string][]*gtd.Task {
	contexts := make(map[string][]*gtd.Task)
	for _, project := range projects {
		if len(project.UnfinishedTasks()) == 0 {
			log.Println("[WARN] Project with no tasks:", project.Name())
		}
		for _, task := range project.Tasks() {
			if !task.Completed {
				if len(task.Contexts) == 0 {
					contexts["default"] = append(contexts["default"], task)
				}
				for _, context := range task.Contexts {
					contexts[context] = append(contexts[context], task)
				}
			}
		}
		external.CreateFile(project.Path(), project.String()) // Persist task checksums to project disk.
	}
	return contexts
}
