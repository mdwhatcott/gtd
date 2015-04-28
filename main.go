package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

var (
	root                     = "/Users/mike/Desktop/GTD"
	runwayLocation           = filepath.Join(root, "0-Runway")
	projectFilesLocation     = filepath.Join(root, "1-Projects")
	processedTasksLocation   = filepath.Join(runwayLocation, "complete-processed")
	completedTasksLocation   = filepath.Join(runwayLocation, "complete")
	stalledProjectsLocation  = filepath.Join(runwayLocation, "projects-stalled")
	finishedProjectsLocation = filepath.Join(runwayLocation, "projects-finished")
)

func main() {
	generateRunway()
	go watchForProjectFileChanges()
	go watchForCompletedTasks()
	waitForInterrupt()
}

func watchForProjectFileChanges() {

}

func watchForCompletedTasks() {
	// poll processedTasksLocation
	// move new files to 'completed'
}

func waitForInterrupt() {
	cancellation := make(chan os.Signal)
	signal.Notify(cancellation, os.Interrupt)
	<-cancellation
	log.Println("Shutting down, generating runway once more...")
	generateRunway()
}

//////////////////////////////////////////////////////////////////////////////

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func fatal(err error) {
	if err != nil {
		_, file, number, _ := runtime.Caller(1)
		log.Printf("Error at %s:%d\n", file, number)
		log.Fatal(err)
	}
}
