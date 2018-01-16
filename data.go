package gtd

import (
	"os"
	"path/filepath"
	"text/template"
)

var (
	RootFolder    = filepath.Join(os.Getenv("HOME"), "Documents/gtd")
	ActionsFolder = filepath.Join(RootFolder, "0-next-actions")
	ProjectFolder = filepath.Join(RootFolder, "1-projects")
	SomedayFolder = filepath.Join(RootFolder, "1-projects-someday")
	MaybeFolder   = filepath.Join(RootFolder, "1-projects-tentative")
)

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
