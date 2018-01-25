package gtd

import (
	"fmt"
	"regexp"
	"strings"
)

type Task struct {
	Text             string
	Completed        bool
	Project          string
	Contexts         []string
	PreviousChecksum string
	CurrentChecksum  string
}

func ParseTask(line string) *Task {
	task := new(Task)
	project, text := extractProjectAndTaskText(line)
	task.Project = project
	task.Text = text
	task.Completed = isCompletedTask(line)
	task.PreviousChecksum = extractPreviousChecksum(line)
	task.Contexts = filterOnPrefix(strings.Fields(line), "@")
	return task
}

func extractProjectAndTaskText(contextLine string) (project, task string) {
	pattern := regexp.MustCompile(`:([a-f0-9]{8}):`)
	results := pattern.Split(contextLine, 5)
	return results[0][len("- [ ]"):], results[1]
}

func extractPreviousChecksum(line string) string {
	pattern := regexp.MustCompile(`:([a-f0-9]{8}):`)
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
	pattern := regexp.MustCompile(` :[a-f0-9]{8}:`)
	index := pattern.FindStringIndex(line)
	if index == nil {
		return line[len("- [ ] "):]
	}
	return line[len("- [ ] "):index[0]]
}
func filterOnPrefix(fields []string, prefix string) (filtered []string) {
	for _, field := range fields {
		if strings.HasPrefix(field, prefix) {
			filtered = append(filtered, strings.ToLower(field))
		}
	}
	return filtered
}

func (this *Task) ProjectString() string {
	builder := new(strings.Builder)
	builder.WriteString("- [")
	if this.Completed {
		builder.WriteString("X")
	} else {
		builder.WriteString(" ")
	}
	builder.WriteString("] ")
	builder.WriteString(this.Text)
	if !this.Completed {
		builder.WriteString(" :")
		builder.WriteString(this.CurrentChecksum)
		builder.WriteString(":")
	}
	return builder.String()
}

func (this *Task) ContextString(projectNameLength int) string {
	format := fmt.Sprintf("- [ ] %%-%ds :%%s: %%s", projectNameLength)
	return fmt.Sprintf(format, this.Project, this.CurrentChecksum, this.Text)
}
