package ux

import (
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

func FormatOutcomeDetails(outcome projections.OutcomeDetails) string {
	prefixes := shortenIDs(actionIDs(outcome.Actions))

	builder := new(strings.Builder)
	builder.WriteString("# ")
	builder.WriteString(outcome.Title)
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("> ")
	builder.WriteString(outcome.Explanation)
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("## Actions:")
	builder.WriteString("\n")
	builder.WriteString("\n")
	for _, action := range outcome.Actions {
		if action.Strategy == core.ActionStrategyConcurrent {
			builder.WriteString("-  ")
		} else {
			builder.WriteString("1. ")
		}
		if action.Status == core.ActionStatusComplete {
			builder.WriteString("[X] ")
		} else if action.Status == core.ActionStatusLatent {
			builder.WriteString("[?] ")
		} else {
			builder.WriteString("[ ] ")
		}
		builder.WriteString("`0x")
		builder.WriteString(prefixes[action.ID])
		builder.WriteString("` ")
		builder.WriteString(action.Description)
		builder.WriteString("\n")
	}
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("## Support Materials:")
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString(outcome.Description)

	return strings.TrimSpace(builder.String())
}

func actionIDs(actions []*projections.ActionDetails) (ids_ []string) {
	for _, action := range actions {
		ids_ = append(ids_, action.ID)
	}
	return ids_
}
