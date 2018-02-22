package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func scanTickler() {
	ticklerMonth := ticklerFolder(time.Now())
	files, err := ioutil.ReadDir(ticklerMonth)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".md") {
			source := filepath.Join(ticklerMonth, name)
			target := filepath.Join(gtd.FolderProjects, name)
			fmt.Println("Tickler project due:", source)
			external.MoveFile(source, target)
		}
	}
}

func ticklerFolder(monthYear time.Time) string {
	return filepath.Join(
		gtd.FolderTickler,
		fmt.Sprint(monthYear.Year()),
		fmt.Sprintf("%02d", int(monthYear.Month())),
	)
}
