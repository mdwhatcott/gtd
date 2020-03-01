package gtd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type Handler interface {
	Handle(interface{})
}

/////////////////////////////////////////////////////////

type Controller struct {
	handler Handler
}

func NewController(handler Handler) *Controller {
	return &Controller{handler: handler}
}

func (this *Controller) ListProjects(input *ListProjectsInputModel) {
	query := new(ListStoredProjectsQuery)
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

type Renderer interface {
	io.Writer
	String(fmt.Stringer)
	JSON(interface{})
}

/////////////////////////////////////////////////////////

type CLIResult struct {
	Stdout io.Writer
	Stderr io.Writer
	Exit   int
}

type RecordingCLIResult struct {
	*CLIResult
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func NewRecordingCLIResult() *RecordingCLIResult {
	out := new(bytes.Buffer)
	err := new(bytes.Buffer)
	return &RecordingCLIResult{
		CLIResult: NewCLIResult(out, err),
		Stdout:    out,
		Stderr:    err,
	}
}

func NewCLIResult(out, err io.Writer) *CLIResult {
	return &CLIResult{Stdout: out, Stderr: err}
}

func (this *CLIResult) String(s fmt.Stringer) {
	_, _ = io.WriteString(this, s.String())
}

func (this *CLIResult) JSON(v interface{}) {
	encoder := json.NewEncoder(this)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}

func (this *CLIResult) Write(v []byte) (int, error) {
	n, err := this.Stdout.Write(v)
	if err != nil {
		this.Exit++
		_, _ = fmt.Fprintf(this.Stderr, "Write err: %v", err)
	}
	return n, err
}

///////////////////////////////////////////////////

type ListStoredProjectsQuery struct {
	Result []StoredProject
}

type StoredProject struct {
	ID      string
	Name    string
	Outcome string
}

////////////////////////////////////////////////////

type TotallyFakeQueryHandler struct{}

func NewTotallyFakeQueryHandler() *TotallyFakeQueryHandler {
	return &TotallyFakeQueryHandler{}
}

func (this *TotallyFakeQueryHandler) Handle(c interface{}) {
	switch context := c.(type) {
	case *ListStoredProjectsQuery:
		context.Result = append(context.Result,
			StoredProject{"1", "name1", "outcome1"},
			StoredProject{"2", "name2", "outcome2"},
			StoredProject{"3", "name3", "outcome3"},
		)
	default:
		panic(fmt.Sprintf("unknown type: %v", reflect.TypeOf(c)))
	}
}
