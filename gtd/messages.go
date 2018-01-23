package gtd

import "fmt"

type CreateProjectCommand struct {
	Blank   bool
	Static  bool
	Name    string
	Outcome string
	Info    string
	Actions []string
}

func (this *CreateProjectCommand) String() string {
	return fmt.Sprintf("%#v", this)
}
func (this *CreateProjectCommand) Set(action string) error {
	this.Actions = append(this.Actions, action)
	return nil
}

type ProjectStatus int

const (
	ProjectActive   ProjectStatus = iota
	ProjectComplete
	ProjectSomeday
	ProjectMaybe
)

type RenegotiateProjectStatusCommand struct {
	ProjectID int
	Result    ProjectStatus
}
