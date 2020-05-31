package ux

import (
	"fmt"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

func FormatIncompleteActions(projection projections.IncompleteActionsByContext) string {
	prefixes := shortenIDs(incompleteActionIDs(projection))

	builder := new(strings.Builder)

	for _, context := range projection.Contexts {
		_, _ = fmt.Fprintf(builder, "## @%s:\n\n", strings.Title(context.Name))

		for _, action := range context.Actions {
			_, _ = fmt.Fprintf(builder,
				"- [ ] `0x%s` %s (`0x%s` %s)\n",
				prefixes[action.ID], action.Description,
				prefixes[action.OutcomeID], action.OutcomeTitle,
			)
		}

		builder.WriteString("\n\n")
	}

	return strings.TrimSpace(builder.String())
}

func incompleteActionIDs(projection projections.IncompleteActionsByContext) (ids_ []string) {
	for _, context := range projection.Contexts {
		for _, action := range context.Actions {
			ids_ = append(ids_, action.OutcomeID, action.ID)
		}
	}
	return ids_
}
