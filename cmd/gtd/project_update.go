package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func updateProjectStatusCLI(inputs []string) {
	flags := external.Flags(usageFlagsUpdateProjects)
	id := flags.Int("id", 0, "The numeric project id (matching output of the last 'gtd project list').")
	status := flags.String("status", "", "The new status of the project (one of: 'complete', 'maybe', 'someday', 'rejected').")
	_ = flags.Parse(inputs)

	if !updateProjectStatus(findProject(*id), *status) {
		exit(flags)
	}
}

func updateProjectStatus(project *gtd.Project, status string) bool {
	if project == nil {
		return false
	} else if destination, found := locations[status]; !found {
		return false
	} else if status != "complete" {
		move(project.Path(), destination)
	} else if recurring := project.RecurringFrequency(); recurring == gtd.RecurringNever {
		move(project.Path(), destination)
	} else {
		prepareForNextOccurrence(project)
	}
	return true
}
func prepareForNextOccurrence(recurring *gtd.Project) {
	for _, task := range recurring.Tasks() {
		task.Completed = false
	}
	source := recurring.Path()
	external.CreateFile(source, recurring.String())
	destination := calculateDestination(recurring.RecurringFrequency())
	move(source, destination)
}
func calculateDestination(recurring gtd.Recurring) string {
	return ticklerFolder(recurring.Next(time.Now()))
}

func move(from, to string) {
	filename := filepath.Base(from)
	destination := filepath.Join(to, filename)
	external.MoveFile(from, destination)
	log.Println("Moved:", destination)
}

func findProject(id int) *gtd.Project {
	projects := gtd.LoadProjects()
	id = id - 1 // use zero-index
	if 0 > id || id >= len(projects) {
		return nil
	}
	return projects[id]
}

var locations = map[string]string{
	"complete": gtd.FolderComplete,
	"someday":  gtd.FolderSomeday,
	"maybe":    gtd.FolderMaybe,
	"reject":   gtd.FolderRejected,
}
