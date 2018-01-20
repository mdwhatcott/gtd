package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func syncTasks() {
	projects := gtd.LoadProjects()
	for _, item := range external.ListDirectory(gtd.FolderActions) {
		scanner := external.ScanFile(filepath.Join(gtd.FolderActions, item.Name()))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			task := gtd.ParseTask(line)
			if task.Completed {
				for _, project := range projects {
					for _, potentialMatch := range project.Tasks() {
						if potentialMatch.PreviousChecksum == task.PreviousChecksum {
							fmt.Println("Crossing off task:", task.Text, task.Project)
							potentialMatch.Completed = true
							external.CreateFile(project.Path(), project.String()) // Persist completed tasks..
						}
					}
				}
			}
		}
	}
}

func sweepTasks() {
	external.DeleteContents(gtd.FolderActions)
	projects := gtd.LoadProjects()
	for context, tasks := range sortTasksByContext(projects) {
		writeTasksInContextToFile(context, tasks)
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
