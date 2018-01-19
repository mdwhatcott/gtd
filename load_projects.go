package gtd

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
)

func LoadProjects() (projects []*Project) {
	dir, err := ioutil.ReadDir(FolderProjects)
	if err != nil {
		log.Fatal(err)
	}
	for i, file := range dir {
		path := filepath.Join(FolderProjects, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Could not read project file:", err)
		}
		projects = append(projects, ParseProject(i+1, path, bytes.NewReader(content)))
	}
	return projects
}
