package main

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mdwhatcott/gtd"
	"github.com/mdwhatcott/gtd/external"
)

func createProjectCLI(inputs []string) {
	createProject(parseCreateProjectCommand(inputs))
}
func parseCreateProjectCommand(inputs []string) (command gtd.CreateProjectCommand) {
	flag := external.Flags(usageFlagsCreateProject)
	flag.BoolVar(&command.Blank, "blank", false, "When set, creates an empty file for the new project.")
	flag.BoolVar(&command.Static, "static", false, "When set, skip editing the new project.")
	flag.StringVar(&command.Name, "name", "", "The succinct, title-case name of the project (use action words).")
	flag.StringVar(&command.Outcome, "outcome", "", "What must become true for this project to be complete?")
	flag.StringVar(&command.Info, "info", "", "What information would you like to record?")
	flag.Var(&command, "action", "What are the next physical action steps to move this project forward?")
	flag.Parse(inputs)
	return command
}

func createProject(command gtd.CreateProjectCommand) {
	path := filepath.Join(gtd.FolderProjects, deriveFilename(command.Name))
	external.CreateFile(path, prepareNewProjectContent(command))

	if !command.Static {
		external.OpenTextEditor(path)
	}
}

func prepareNewProjectContent(command gtd.CreateProjectCommand) string {
	if command.Blank {
		return ""
	} else {
		return external.ExecuteTemplate(gtd.ProjectTemplate, command)
	}
}

func deriveFilename(name string) string {
	name = replace(name, "  ", " ")
	name = strings.TrimSpace(name)
	name = strings.Map(toLowerASCII, name)
	name = replace(name, "--", "-")
	name = strings.Trim(name, "-")
	return name + ".md"
}
func toLowerASCII(r rune) rune {
	if unicode.IsSpace(r) {
		return '-'
	} else if r < unicode.MaxLatin1 && unicode.IsLetter(r) {
		return unicode.ToLower(r)
	} else {
		return -1 // discard
	}
}
func replace(input, old, new string) string {
	for strings.Contains(input, old) {
		input = strings.Replace(input, old, new, -1)
	}
	return input
}
