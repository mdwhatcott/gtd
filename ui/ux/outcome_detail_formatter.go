package ux

import (
	"strings"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func FormatOutcomeDetails(outcome projections.OutcomeDetails) string {
	PREFIXES := shortenIDs(actionIDs(outcome.Actions))

	BUILDER := new(strings.Builder)
	BUILDER.WriteString("# ")
	BUILDER.WriteString(outcome.Title)
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("> ")
	BUILDER.WriteString(outcome.Explanation)
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("## Actions:")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	for _, ACTION := range outcome.Actions {
		if ACTION.Strategy == core.ActionStrategyConcurrent {
			BUILDER.WriteString("-  ")
		} else {
			BUILDER.WriteString("1. ")
		}
		if ACTION.Status == core.ActionStatusComplete {
			BUILDER.WriteString("[X] ")
		} else if ACTION.Status == core.ActionStatusLatent {
			BUILDER.WriteString("[?] ")
		} else {
			BUILDER.WriteString("[ ] ")
		}
		BUILDER.WriteString("`0x")
		BUILDER.WriteString(PREFIXES[ACTION.ID])
		BUILDER.WriteString("` ")
		BUILDER.WriteString(ACTION.Description)
		BUILDER.WriteString("\n")
	}
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("## Support Materials:")
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	BUILDER.WriteString(outcome.Description)

	return strings.TrimSpace(BUILDER.String())
}

func actionIDs(actions []*projections.ActionDetails) (ids_ []string) {
	for _, ACTION := range actions {
		ids_ = append(ids_, ACTION.ID)
	}
	return ids_
}
