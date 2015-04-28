package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mdwhatcott/gtd/projects"
	"github.com/mdwhatcott/gtd/runway"
)

var (
	root                   = "/Users/mike/Dropbox/GTD"
	runwayLocation         = filepath.Join(root, "0-Runway")
	projectFilesLocation   = filepath.Join(root, "1-Projects")
	processedTasksLocation = filepath.Join(runwayLocation, "complete-processed")
)

func main() {
	generateRunway(runwayLocation, projectFilesLocation, processedTasksLocation)
	go watchForProjectFileChanges(projectFilesLocation)
	watchForCompletedTasks(processedTasksLocation, projectFilesLocation)
}

func generateRunway(runwayLocation, projectFilesLocation, processedTasksLocation string) {
	log.Printf("Generating the runway at: %s\n", runwayLocation)
	fatal(os.RemoveAll(runwayLocation))
	makeDirectory(runwayLocation)
	makeDirectory(processedTasksLocation)

	var (
		completedTasksLocation   = makeDirectory(filepath.Join(runwayLocation, "complete"))
		stalledProjectsLocation  = makeDirectory(filepath.Join(runwayLocation, "projects-stalled"))
		finishedProjectsLocation = makeDirectory(filepath.Join(runwayLocation, "projects-finished"))
	)

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

	for context, listing := range runway.GroupContextListings(projectListing) {
		location := makeDirectory(filepath.Join(runwayLocation, context))
		for _, task := range listing {
			createFile(location, taskFilename(task))
		}
	}
	for _, task := range runway.IdentifyCompletedTasks(projectListing) {
		createFile(completedTasksLocation, taskFilename(task))
	}
	for _, stalled := range runway.IdentifyStalledProjects(projectListing) {
		createFile(stalledProjectsLocation, projectFilename(stalled))
	}
	for _, finished := range runway.IdentifyFinishedProjects(projectListing) {
		createFile(finishedProjectsLocation, finished.Name)
	}
}

func watchForProjectFileChanges(projectsFolder string) {

}

func watchForCompletedTasks(processedTasksLocation, projectFilesLocation string) {

}

//////////////////////////////////////////////////////////////////////////////

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func isDir(folder string) bool {
	info, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fatal(err error) {
	if err != nil {
		_, file, number, _ := runtime.Caller(1)
		log.Printf("Error at %s:%d\n", file, number)
		log.Fatal(err)
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
func projectFilename(project projects.Project) string {
	return fmt.Sprintf("%s")
}
