package main

import (
	"fmt"
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
	flags.Parse(inputs)

	project := findProject(*id)
	if !updateProjectStatus(project, *status) {
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
		for _, task := range project.Tasks() {
			task.Completed = false
		}
		external.CreateFile(project.Path(), project.String())
		move(project.Path(), calculateDestination(recurring))
	}
	return true
}
func calculateDestination(recurring gtd.Recurring) string {
	next := recurring.Next(time.Now())
	return filepath.Join(
		gtd.FolderTickler,
		fmt.Sprint(next.Year()),
		fmt.Sprintf("%02d", int(next.Month())),
	)
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
