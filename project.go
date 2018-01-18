package gtd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	id    int
	path  string
	name  string
	tasks []*Task
}

func ParseProject(id int, path string, reader io.Reader) *Project {
	project := &Project{
		id:   id,
		path: path,
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") && project.name == "" {
			project.name = line[2:]
		}

		if isTask(line) {
			project.tasks = append(project.tasks, &Task{
				Text:             taskText(line),
				Completed:        isCompletedTask(line),
				Project:          project.name,
				Contexts:         filterOnPrefix(strings.Fields(line), "@"),
				previousChecksum: extractPreviousChecksum(line),
				currentChecksum:  simpleHash(line),
			})
		}
	}
	return project
}
func extractPreviousChecksum(line string) string {
	pattern := regexp.MustCompile(`{{(.+)}}`)
	matches := pattern.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
func isTask(line string) bool {
	matched, _ := regexp.MatchString(`^- \[[xX ]] .+`, line)
	return matched
}
func isCompletedTask(line string) bool {
	matched, _ := regexp.MatchString(`^- \[x|X] .+`, line)
	return matched
}
func taskText(line string) string {
	taskPrefix := "- [ ] "
	checksumIndex := strings.Index(line, "{{")
	if checksumIndex < 0 {
		checksumIndex = len(line)
	}
	return strings.TrimSpace(line[len(taskPrefix):checksumIndex])
}
func filterOnPrefix(fields []string, prefix string) (filtered []string) {
	for _, field := range fields {
		if strings.HasPrefix(field, prefix) {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func (p *Project) Path() string {
	return p.path
}

func (this *Project) Name() string {
	return fmt.Sprintf(listFormat, this.id, this.listingName())
}

func (this *Project) listingName() string {
	if this.name != "" {
		return this.name
	} else {
		file := filepath.Base(this.path)
		return strings.TrimRight(file, filepath.Ext(file))
	}
}

const listFormat = "%-4d%s"

func (this *Project) Tasks() []*Task {
	return this.tasks
}

type Task struct {
	Text             string
	Completed        bool
	Project          string
	Contexts         []string
	previousChecksum string
	currentChecksum  string
}

func simpleHash(input string) string {
	hash := fnv.New32a()
	hash.Write([]byte(input))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
