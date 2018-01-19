package main

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func sweepTasksCLI(input []string) {
	external.Flags(usageFlagsSweepTasks).Parse(input)
	sweepTasks(gtd.LoadProjects())
}

func sweepTasks(projects []*gtd.Project) {
	for context, tasks := range sortTasksByContext(projects) {
		writeTasksInContextToFile(context, tasks)
	}
}

func sortTasksByContext(projects []*gtd.Project) map[string][]*gtd.Task {
	contexts := make(map[string][]*gtd.Task)
	for _, project := range projects {
		tasks := project.Tasks()
		if len(tasks) == 0 {
			log.Println("[WARN] Project with no tasks:", project.Name())
		}
		for _, task := range tasks {
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
func writeTasksInContextToFile(context string, tasks []*gtd.Task) {
	maxProjectNameLength := 0
	for _, task := range tasks {
		if len(task.Project) > maxProjectNameLength {
			maxProjectNameLength = len(task.Project)
		}
	}
	buffer := new(bytes.Buffer)
	for _, task := range tasks {
		buffer.WriteString(task.ContextString(maxProjectNameLength) + "\n")
	}

	path := filepath.Join(gtd.FolderActions, strings.ToLower(strings.Trim(context, "@"))+".md")
	external.CreateFile(path, buffer.String())
}
