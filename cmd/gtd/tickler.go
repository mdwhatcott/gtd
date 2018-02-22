package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

func scanTickler() {
	now := time.Now()
	err := filepath.Walk(gtd.FolderTickler, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		if project := ParseTicklerProject(path); project.IsDue(now) {
			log.Println("Project from tickler has become due (moving to active projects folder):", info.Name())
			external.MoveFile(path, project.TargetPath())
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

type TicklerProject struct {
	path string
	due  time.Time
}

func ParseTicklerProject(path string) TicklerProject {
	parts := strings.Split(path, string(filepath.Separator))
	month, _ := strconv.Atoi(parts[len(parts)-2])
	year, _ := strconv.Atoi(parts[len(parts)-3])
	monthYear := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	return TicklerProject{
		path: path,
		due:  monthYear,
	}
}

func (this TicklerProject) IsDue(now time.Time) bool {
	return this.due.Month() == now.Month() && this.due.Year() == now.Year()
}

func (this TicklerProject) TargetPath() string {
	return filepath.Join(gtd.FolderProjects, filepath.Base(this.path))
}
