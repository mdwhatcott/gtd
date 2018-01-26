package storage

type ContextLoader interface {
	AppendToContexts(task *Task)
	FinishedTasks() ([]*Task, error)
	Save(context *Context) error
}
