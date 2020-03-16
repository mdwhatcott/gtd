package gtd

import (
	"fmt"
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/storage"
)

type Controller struct {
	handler Handler
}

func NewController(handler Handler) *Controller {
	return &Controller{handler: handler}
}

func (this *Controller) ListProjects(input *ListProjectsInputModel) {
	query := new(storage.ListStoredProjectsQuery)
	this.handler.Handle(query)
	for _, project := range query.Result {
		input.Results.Projects = append(input.Results.Projects, ProjectForCLI{
			ID:      project.ID,
			Name:    project.Name,
			Outcome: project.Outcome,
		})
	}
}

/////////////////////////////////////////////////////////

type ListProjectsInputModel struct {
	Results struct {
		Projects []ProjectForCLI
	}
}

func NewListProjectsInputModel() *ListProjectsInputModel {
	return &ListProjectsInputModel{}
}

func (this *ListProjectsInputModel) Render(result Renderer) {
	for _, project := range this.Results.Projects {
		result.String(project)
	}
}

type ProjectForCLI struct {
	ID      string
	Name    string
	Outcome string
}

func (this ProjectForCLI) String() string {
	return fmt.Sprintf("%s. %s\n   - %s\n", this.ID, this.Name, this.Outcome)
}

/////////////////////////////////////////////////////////

type TotallyFakeQueryHandler struct{}

func NewTotallyFakeQueryHandler() *TotallyFakeQueryHandler {
	return &TotallyFakeQueryHandler{}
}

func (this *TotallyFakeQueryHandler) Handle(c interface{}) {
	switch context := c.(type) {
	case *storage.ListStoredProjectsQuery:
		context.Result = append(context.Result,
			storage.StoredProject{ID: "1", Name: "name1", Outcome: "outcome1"},
			storage.StoredProject{ID: "2", Name: "name2", Outcome: "outcome2"},
			storage.StoredProject{ID: "3", Name: "name3", Outcome: "outcome3"},
		)
	default:
		panic(fmt.Sprintf("unknown type: %v", reflect.TypeOf(c)))
	}
}
