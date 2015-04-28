package projects

import (
	"bufio"
	"strings"
)

func ParseTasks(project, content string) (tasks []Task) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := trim(scanner.Text())
		if strings.HasPrefix(line, "- [ ] ") {
			tasks = append(tasks, incompleteTask(project, line, len(tasks)))
		} else if strings.HasPrefix(line, "- [X] ") {
			tasks = append(tasks, completeTask(project, line, len(tasks)))
		}
	}
	return tasks
}

func trim(line string) string {
	for strings.Contains(line, "  ") {
		line = strings.Replace(line, "  ", " ", -1)
	}
	return strings.TrimSpace(line)
}

func completeTask(project, line string, index int) Task {
	task := incompleteTask(project, line, index)
	task.Complete = true
	return task
}

func incompleteTask(project, line string, index int) Task {
	task := Task{
		ParentProject: project,
		Index:         index,
		Text:          line[len("- [ ] "):],
		Contexts:      []string{},
	}
	for _, word := range strings.Split(task.Text, " ") {
		if strings.HasPrefix(word, "@") {
			task.Contexts = append(task.Contexts, strings.ToLower(word))
		}
	}
	return task
}
