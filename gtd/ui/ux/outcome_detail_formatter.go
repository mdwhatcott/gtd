package ux

import (
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

func FormatOutcomeDetails(outcome *projections.OutcomeDetails) string {
	//builder := new(strings.Builder)
	//builder.WriteString("# ")
	//builder.WriteString(outcome.Title)
	//builder.WriteString("\n")
	//TODO
	//return builder.String()
	return strings.Join([]string{
		"# Title",
		"",
		"> Explanation",
		"",
		"",
		"## Actions:",
		"",
		"-  [ ] `0x0001` Action 1",
		"-  [?] `0x0002` Action 2",
		"-  [X] `0x0003` Action 3",
		"1. [ ] `0x0004` Action 4",
		"1. [?] `0x0005` Action 5",
		"1. [X] `0x0006` Action 6",
	}, "\n")
}
