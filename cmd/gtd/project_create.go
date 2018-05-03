package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/mdwhatcott/gtd/external"
	"github.com/mdwhatcott/gtd/gtd"
)

type NewProjectInfo struct {
	Blank   bool
	Static  bool
	Name    string
	Outcome string
	Info    string
	Actions []string
}

func (this *NewProjectInfo) String() string {
	return fmt.Sprintf("%#v", this)
}
func (this *NewProjectInfo) Set(action string) error {
	this.Actions = append(this.Actions, action)
	return nil
}

func createProjects() {
	for {
		fmt.Print("Enter new project name (<blank> to quit): ")
		name := external.ReadLine()
		if name == "" {
			break
		}
		createProject(NewProjectInfo{Name: name})
	}
}

func createProject(command NewProjectInfo) {
	path := filepath.Join(gtd.FolderProjects, deriveFilename(command.Name))
	external.CreateFile(path, prepareNewProjectContent(command))

	if !command.Static {
		external.OpenTextEditorAndWait(path)
	}
}

func prepareNewProjectContent(command NewProjectInfo) string {
	if command.Blank {
		return ""
	} else {
		return external.ExecuteTemplate(ProjectTemplate, command)
	}
}

func deriveFilename(name string) string {
	name = replace(name, "  ", " ")
	name = strings.TrimSpace(name)
	name = replace(name, ":", "-")
	name = strings.Map(toLowerASCII, name)
	name = strings.Trim(name, "-")
	return name + ".md"
}
func replace(input, old, new string) string {
	for strings.Contains(input, old) {
		input = strings.Replace(input, old, new, -1)
	}
	return input
}
func toLowerASCII(r rune) rune {
	switch {
	case isDigit(r):
		return r
	case separators[r]:
		return '-'
	case isLowercase(r):
		return r
	case isUppercase(r):
		return unicode.ToLower(r)
	default:
		return discardRune
	}
}

func isDigit(r rune) bool     { return '0' <= r && r <= '9' }
func isLowercase(r rune) bool { return 'a' <= r && r <= 'z' }
func isUppercase(r rune) bool { return 'A' <= r && r <= 'Z' }

var separators = map[rune]bool{
	'-':  true,
	'_':  true,
	' ':  true,
	'\t': true,
	'\n': true,
}

const discardRune = -1

var ProjectTemplate = template.Must(template.New("project").Parse(`# {{.Name}}

Desired Outcome: {{.Outcome}}

RECURRING: ???

## Info

{{.Info}}


## Tasks:

{{range .Actions}}
- [ ] {{.}}{{end}}


## Activity Log:


### 2018-01-15

- What happened today?`))
