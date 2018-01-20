package gtd

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"path/filepath"
	"strings"
)

type Project struct {
	id    int
	path  string
	name  string
	tasks []*Task
	lines []string
}

func ParseProject(id int, path string, reader io.Reader) *Project {
	project := &Project{
		id:   id,
		path: path,
	}

	for scanner := bufio.NewScanner(reader); scanner.Scan(); {
		line := scanner.Text()
		project.lines = append(project.lines, line)

		if strings.HasPrefix(line, "# ") && project.name == "" {
			project.name = line[2:]
		}

		if isTask(line) {
			project.tasks = append(project.tasks, &Task{
				Text:             taskText(line),
				Completed:        isCompletedTask(line),
				Project:          project.name,
				Contexts:         filterOnPrefix(strings.Fields(line), "@"),
				PreviousChecksum: extractPreviousChecksum(line),
				CurrentChecksum:  checksum(project.name + " " + taskText(line)),
			})
		}
	}

	return project
}

func checksum(input string) string {
	hash := fnv.New32a()
	hash.Write([]byte(input))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

func (p *Project) Path() string {
	return p.path
}
func (this *Project) Name() string {
	const listFormat = "%-4d%s"
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

func (this *Project) Tasks() []*Task {
	return this.tasks
}

func (this *Project) UnfinishedTasks() (unfinished []*Task) {
	for _, task := range this.Tasks() {
		if !task.Completed {
			unfinished = append(unfinished, task)
		}
	}
	return unfinished
}

func (this *Project) String() string {
	writer := new(bytes.Buffer)
	taskIndex := 0
	for _, line := range this.lines {
		if isTask(line) {
			writer.WriteString(this.tasks[taskIndex].ProjectString())
			taskIndex++
		} else {
			writer.WriteString(line)
		}
		writer.WriteString("\n")
	}
	return writer.String()
}
