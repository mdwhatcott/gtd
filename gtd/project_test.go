package gtd

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestProjectFixture(t *testing.T) {
	gunit.Run(new(ProjectFixture), t)
}

type ProjectFixture struct {
	*gunit.Fixture
}

func (this *ProjectFixture) TestParseProjectName() {
	content := strings.NewReader("# I am a project name\n\nWith important stuff\n# not a project name")
	project := ParseProject(42, "path", content)
	this.So(project.Name(), should.Equal, "42  I am a project name")
}

func (this *ProjectFixture) TestParseProjectNoName() {
	content := strings.NewReader("I am a project name but I'm not marked as such (#)")
	project := ParseProject(42, "/somewhere/path.md", content)
	this.So(project.Name(), should.Equal, "42  path")
}

func (this *ProjectFixture) TestParseTasks() {
	content := strings.NewReader(`
# Title

Info

## Some tasks

- [X] finished 1 :a1b2c3d4:
- [ ] unfinished 1 @HomeDepot

- [ ] unfinished 2
- [x] finished 2 @@Person @Phone

I'm not - [ ] a task
`)
	project := ParseProject(42, "path", content)
	this.So(project.Tasks(), should.Resemble, []*Task{
		{
			Text:             "finished 1",
			Completed:        true,
			Project:          "Title",
			PreviousChecksum: "a1b2c3d4",
			CurrentChecksum:  "3a63954a",
		},
		{
			Text:            "unfinished 1 @HomeDepot",
			Completed:       false,
			Project:         "Title",
			Contexts:        []string{"@HomeDepot"},
			CurrentChecksum: "eae9c550",
		},
		{
			Text:            "unfinished 2",
			Completed:       false,
			Project:         "Title",
			CurrentChecksum: "9e5fc0a4",
		},
		{
			Text:            "finished 2 @@Person @Phone",
			Completed:       true,
			Project:         "Title",
			Contexts:        []string{"@@Person", "@Phone"},
			CurrentChecksum: "6b9ebeb4",
		},
	})
}
