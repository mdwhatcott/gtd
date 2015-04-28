package projects

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestTaskParsing(t *testing.T) {
	var cases = []ParseCase{
		{
			description: "several tasks",
			input: strings.Join([]string{
				"blah blah blah not a task",
				"- [ ] 0",
				"Not a task",
				"- [X] 1",
				"- Also not a task",
				"- [ ] 2 @A",
				"- [ Almost a task",
				"- [X] 3     @B    ",
				"- [ [ So close",
				"- [X] 4 @A @B",
				"- [X]No space, no cigar",
				"- [ ] @C 5",
			}, "\n"),
			expected: []Task{
				Task{
					Text:          "0",
					Complete:      false,
					Contexts:      []string{},
					ParentProject: "project-name",
					Index:         0,
				},
				Task{
					Text:          "1",
					Complete:      true,
					Contexts:      []string{},
					ParentProject: "project-name",
					Index:         1,
				},
				Task{
					Text:          "2 @A",
					Complete:      false,
					Contexts:      []string{"@a"},
					ParentProject: "project-name",
					Index:         2,
				},
				Task{
					Text:          "3 @B",
					Complete:      true,
					Contexts:      []string{"@b"},
					ParentProject: "project-name",
					Index:         3,
				},
				Task{
					Text:          "4 @A @B",
					Complete:      true,
					Contexts:      []string{"@a", "@b"},
					ParentProject: "project-name",
					Index:         4,
				},
				Task{
					Text:          "@C 5",
					Complete:      false,
					Contexts:      []string{"@c"},
					ParentProject: "project-name",
					Index:         5,
				},
			},
		},
	}

	for i, test := range cases {
		tasks := ParseTasks("project-name", test.input)
		if ok, message := assertions.So(tasks, should.Resemble, test.expected); !ok {
			t.Errorf("Test case #%d:\n%s", i, message)
		}
	}
}

type ParseCase struct {
	description string
	input       string
	expected    []Task
}
