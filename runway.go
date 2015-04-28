package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/gtd/projects"
	"github.com/mdwhatcott/gtd/runway"
)

func generateRunway() {
	log.Printf("Generating the runway at: %s\n", runwayLocation)
	prepareWorkspace()
	projectListing := parseAllProjects()
	createContextListings(projectListing)
	createCompletedTaskListing(projectListing)
	createStalledProjectListing(projectListing)
	createFinishedProjectListing(projectListing)
}

func prepareWorkspace() {
	fatal(os.RemoveAll(runwayLocation))
	makeDirectory(runwayLocation)
	makeDirectory(processedTasksLocation)
	makeDirectory(completedTasksLocation)
	makeDirectory(stalledProjectsLocation)
	makeDirectory(finishedProjectsLocation)
}
func parseAllProjects() []projects.Project {
	listing, err := ioutil.ReadDir(projectFilesLocation)
	fatal(err)
	projectListing := []projects.Project{}
	for _, info := range listing {
		content, err := ioutil.ReadFile(filepath.Join(projectFilesLocation, info.Name()))
		fatal(err)
		name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		projectListing = append(projectListing, projects.Project{
			Name:  name,
			Tasks: projects.ParseTasks(name, string(content)),
		})
	}
	return projectListing
}
func createContextListings(projectListing []projects.Project) {
	for context, listing := range runway.GroupContextListings(projectListing) {
		location := makeDirectory(filepath.Join(runwayLocation, context))
		for _, task := range listing {
			createFile(location, taskFilename(task))
		}
	}
}
func createCompletedTaskListing(projectListing []projects.Project) {
	for _, task := range runway.IdentifyCompletedTasks(projectListing) {
		createFile(completedTasksLocation, taskFilename(task))
	}
}
func createStalledProjectListing(projectListing []projects.Project) {
	for _, stalled := range runway.IdentifyStalledProjects(projectListing) {
		createFile(stalledProjectsLocation, stalled.Name)
	}
}
func createFinishedProjectListing(projectListing []projects.Project) {
	for _, finished := range runway.IdentifyFinishedProjects(projectListing) {
		createFile(finishedProjectsLocation, finished.Name)
	}
}

func makeDirectory(location string) string {
	fatal(os.Mkdir(location, 0644))
	return location
}

func createFile(location, filename string) {
	file, err := os.Create(filepath.Join(filename))
	fatal(err)
	defer file.Close()
}

func taskFilename(task projects.Task) string {
	return fmt.Sprintf("%s_%d_%s", task.ParentProject, task.Index, task.Text) // TODO: remove unsafe characters
}
