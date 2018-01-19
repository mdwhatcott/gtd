package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func syncTasksCLI(input []string) {
	external.Flags(usageFlagsSyncTasks).Parse(input)
	syncTasks(gtd.LoadProjects())
}

func syncTasks(projects []*gtd.Project) {
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
