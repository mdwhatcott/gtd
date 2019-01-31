package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func tasksCLI(input []string) {
	flag := external.Flags(usageFlagsTasks)
	_ = flag.Parse(input)

	first, _ := firstAndRemaining(flag.Args())

	switch first {
	case "":
		regenerateTasks()
	case "review":
		reviewTasks()
	default:
		exit(flag)
	}
}
func regenerateTasks() {
	syncTasks()
	sweepTasks()
}

func syncTasks() {
	projects := gtd.LoadProjects()
	for _, item := range external.ListDirectory(gtd.FolderActions) {
		scanner := external.ScanFile(join(gtd.FolderActions, item.Name()))
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
							log.Printf("Crossing off: %s (from: %s)\n", task.Text, strings.TrimSpace(task.Project))
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
	builder := new(strings.Builder)
	for _, task := range tasks {
		builder.WriteString(task.ContextString(maxProjectNameLength) + "\n")
	}

	contextName := strings.Trim(context, "@")
	contextName = strings.ToLower(contextName)
	filename := fmt.Sprintf("%s.md", contextName)
	path := join(gtd.FolderActions, filename)
	external.CreateFile(path, builder.String())
}

func reviewTasks() {
	files, err := ioutil.ReadDir(gtd.FolderActions)
	if err != nil {
		log.Fatalln("Could not read actions folder.")
	}

	for _, info := range files {
		if !strings.HasSuffix(info.Name(), ".md") {
			continue
		}

		fmt.Println("Now reviewing", info.Name())
		external.OpenTextEditorAndWait(join(gtd.FolderActions, info.Name()))
	}
}
