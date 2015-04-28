package runway

import (
	"testing"

	"github.com/mdwhatcott/gtd/projects"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

var (
	hi_blah_incomplete = projects.Task{
		Text:          "hi there @blah",
		Contexts:      []string{"@blah"},
		ParentProject: "A",
		Complete:      false,
		Index:         0,
	}
	hi_there_blah_complete = projects.Task{
		Text:          "hi there @blah",
		Contexts:      []string{"@blah"},
		ParentProject: "A",
		Complete:      true,
		Index:         1,
	}
	bye_foo_blah_incomplete = projects.Task{
		Text:          "bye @foo @blah",
		Contexts:      []string{"@foo", "@blah"},
		ParentProject: "A",
		Complete:      false,
		Index:         2,
	}
	bye_foo_blah_complete = projects.Task{
		Text:          "bye @foo @blah",
		Contexts:      []string{"@foo", "@blah"},
		ParentProject: "A",
		Complete:      true,
		Index:         3,
	}

	hi_incomplete = projects.Task{
		Text:          "hi",
		Contexts:      []string{},
		ParentProject: "Stalled",
		Complete:      false,
		Index:         0,
	}

	hi_complete = projects.Task{
		Text:          "hi",
		Contexts:      []string{},
		ParentProject: "Finished",
		Complete:      true,
		Index:         0,
	}

	project_A = projects.Project{
		Name: "A",
		Tasks: []projects.Task{
			hi_blah_incomplete,
			hi_there_blah_complete,
			bye_foo_blah_incomplete,
			bye_foo_blah_complete,
		},
	}
	project_stalled = projects.Project{
		Name:  "Stalled",
		Tasks: []projects.Task{hi_incomplete},
	}
	project_finished = projects.Project{
		Name:  "Finished",
		Tasks: []projects.Task{hi_complete},
	}

	input = []projects.Project{
		project_A,
		project_stalled,
		project_finished,
	}
)

func TestGroupContextListings(t *testing.T) {
	expected := map[string][]projects.Task{
		"@blah": []projects.Task{hi_blah_incomplete, bye_foo_blah_incomplete},
		"@foo":  []projects.Task{bye_foo_blah_incomplete},
	}
	actual := GroupContextListings(input)
	if ok, message := assertions.So(actual, should.Resemble, expected); !ok {
		t.Error("\n" + message)
	}
}

func TestIdentifyCompletedTasks(t *testing.T) {
	expected := []projects.Task{
		hi_there_blah_complete,
		bye_foo_blah_complete,
		hi_complete,
	}
	actual := IdentifyCompletedTasks(input)
	if ok, message := assertions.So(actual, should.Resemble, expected); !ok {
		t.Error("\n" + message)
	}
}

func TestIdentifyStalledProjects(t *testing.T) {
	expected := []projects.Project{project_stalled}
	actual := IdentifyStalledProjects(input)
	if ok, message := assertions.So(actual, should.Resemble, expected); !ok {
		t.Error("\n" + message)
	}
}

func TestIdentifyFinishedProjects(t *testing.T) {
	expected := []projects.Project{project_finished}
	actual := IdentifyFinishedProjects(input)
	if ok, message := assertions.So(actual, should.Resemble, expected); !ok {
		t.Error("\n" + message)
	}
}
