package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mdwhatcott/gtd/projects"
)

var running bool = true

var (
	root                     = "/Users/mike/Desktop/GTD"
	runwayLocation           = filepath.Join(root, "0-Runway")
	projectFilesLocation     = filepath.Join(root, "1-Projects")
	processedTasksLocation   = filepath.Join(runwayLocation, "complete-processed")
	completedTasksLocation   = filepath.Join(runwayLocation, "complete")
	stalledProjectsLocation  = filepath.Join(runwayLocation, "projects-stalled")
	finishedProjectsLocation = filepath.Join(runwayLocation, "projects-finished")
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	generateRunway()
	go watchForProjectFileChanges()
	go watchForCompletedTasks()
	waitForInterrupt()
}

func watchForProjectFileChanges() {
	checksummer := projects.NewChecksummer()
	for running {
		listing, err := ioutil.ReadDir(projectFilesLocation)
		fatal(err)
		if checksummer.IsDirty(listing) {
			generateRunway()
		}
		time.Sleep(time.Second)
	}
}

func watchForCompletedTasks() {
	for running {
		listing, err := ioutil.ReadDir(completedTasksLocation)
		fatal(err)
		for _, completed := range listing {
			completeTask(completed.Name())
		}
		time.Sleep(time.Second)
	}
}

func completeTask(filename string) {
	taskPath := filepath.Join(completedTasksLocation, filename)
	completedPath := filepath.Join(processedTasksLocation, filename)
	fatal(os.Rename(taskPath, completedPath))

	fields := strings.Split(filename, "__")
	projectName := fields[0]
	taskIndex, err := strconv.Atoi(fields[1])
	fatal(err)
	projectPath := filepath.Join(projectFilesLocation, projectName+".md")
	projectBytes, err := ioutil.ReadFile(projectPath)
	fatal(err)
	reader := strings.NewReader(string(projectBytes))
	writer := new(bytes.Buffer)
	scanner := bufio.NewScanner(reader)
	tasks := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "- [ ] ") || strings.HasPrefix(line, "- [X] ") {
			if tasks == taskIndex {
				line = strings.TrimSpace(strings.Replace(line, "- [ ] ", "- [X] ", 1))
			}
			tasks += 1
		}
		writer.WriteString(strings.TrimSpace(line) + "\n")
	}
	fatal(ioutil.WriteFile(projectPath, writer.Bytes(), 0644))
}

func waitForInterrupt() {
	cancellation := make(chan os.Signal)
	signal.Notify(cancellation, os.Interrupt)
	<-cancellation
	log.Println("Shutting down, generating runway once more...")
	running = false
	time.Sleep(time.Second)
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
