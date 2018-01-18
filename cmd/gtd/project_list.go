package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/mdwhatcott/gtd"
)

func listProjects(inputs []string) {
	flag := flags(usageFlagsListProjects)
	review := flag.Bool("review", false, "When set, review each project via a REPL and text editor sessions.")
	flag.Parse(inputs)

	projects := LoadProjects()
	for _, project := range projects {
		fmt.Println(project.Name())
		if *review {
			edit(project.Path())
			// TODO: would you like to update the status of the project? (complete, defer, etc...)
		}
	}
}

// TODO: move to top-level lib (mock out fs?)
func LoadProjects() (projects []*gtd.Project) {
	dir, err := ioutil.ReadDir(gtd.FolderProjects)
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range dir {
		path := filepath.Join(gtd.FolderProjects, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Could not read project file:", err)
		}
		projects = append(projects, gtd.ParseProject(i+1, path, bytes.NewReader(content)))
	}
	return projects
}
