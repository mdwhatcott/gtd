package storage

type Project struct {
	id        int
	path      string
	name      string
	tasks     []*Task
	lines     []string
	recurring Recurring
}
