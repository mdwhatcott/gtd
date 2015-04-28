package projects

type Project struct {
	Name  string
	Tasks []Task
}

type Task struct {
	Text          string
	Complete      bool
	Contexts      []string
	ParentProject string
	Index         int
}
