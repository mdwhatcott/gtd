package gtd

import (
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"
)

var (
	Now            = time.Now()
	FolderRoot     = root(os.Getenv("GTDPATH"))
	FolderActions  = filepath.Join(FolderRoot, "0-next-actions")
	FolderProjects = filepath.Join(FolderRoot, "1-projects")
	FolderArchive  = filepath.Join(FolderRoot, "1-projects-archive", strconv.Itoa(Now.Year()))
	FolderSomeday  = filepath.Join(FolderRoot, "1-projects-someday")
	FolderTickler  = filepath.Join(FolderSomeday, "1-tickler")
	FolderMaybe    = filepath.Join(FolderRoot, "1-projects-tentative")
)

func root(proposed string) string {
	if proposed != "" {
		return proposed
	}
	return filepath.Join(os.Getenv("HOME"), "Documents", "gtd")
}

var ProjectTemplate = template.Must(template.New("project").Parse(`# {{.Name}}

Desired Outcome: {{.Outcome}}


## Info

{{.Info}}


## Tasks:

{{range .Actions}}
- [ ] {{.}}{{end}}


## Activity Log:


### 2018-01-15

- What happened today?`))
