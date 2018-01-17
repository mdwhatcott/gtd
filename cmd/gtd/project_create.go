package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/gtd"
)

func createProject(inputs []string) {
	command := parseCreateProjectCommand(inputs)
	path := filepath.Join(gtd.FolderProjects, deriveFilename(command))

	if err := ioutil.WriteFile(path, prepareNewProjectContent(command), 0644); err != nil {
		log.Fatalln("Could not write project file:", err)
	}

	if !command.Static {
		if err := exec.Command("subl", "--wait", path).Run(); err != nil {
			log.Fatalln("Error editing project file:", err)
		}
	}
}

func parseCreateProjectCommand(inputs []string) (command gtd.CreateProjectCommand) {
	createProjectFlags.BoolVar(&command.Blank, "blank", false, "When set, creates an empty file for the new project.")
	createProjectFlags.BoolVar(&command.Static, "static", false, "When set, skip editing the new project.")
	createProjectFlags.StringVar(&command.Name, "name", "", "The succinct, title-case name of the project (use action words).")
	createProjectFlags.StringVar(&command.Outcome, "outcome", "", "What must become true for this project to be complete?")
	createProjectFlags.StringVar(&command.Info, "info", "", "What information would you like to record?")
	createProjectFlags.Var(&command, "action", "What are the next physical action steps to move this project forward?")
	createProjectFlags.Parse(inputs)
	return command
}

func prepareNewProjectContent(command gtd.CreateProjectCommand) []byte {
	if command.Blank {
		return nil
	} else {
		var content bytes.Buffer
		if err := gtd.ProjectTemplate.Execute(&content, command); err != nil {
			log.Fatalln("Could not execute template:", err)
		}
		return content.Bytes()
	}
}

func deriveFilename(command gtd.CreateProjectCommand) string {
	name := command.Name
	name = filenameReplacer.Replace(name)
	name = strings.TrimSpace(name)
	name = replace(name, " ", "-")
	name = strings.ToLower(name)
	return name + ".md"
}

func replace(input, old, new string) string {
	return strings.Replace(input, old, new, -1)
}

var filenameReplacer = strings.NewReplacer(
	"/", " ",
	`|`, " ",
	":", " ",
	"?", " ",
	"%", " ",
	"*", " ",
	"|", " ",
	`"`, " ",
	"<", " ",
	">", " ",
	".", " ",
	"  ", " ",
)
