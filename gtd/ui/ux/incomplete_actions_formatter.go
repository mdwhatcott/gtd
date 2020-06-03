package ux

import (
	"fmt"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

func FormatIncompleteActions(contexts ...*projections.Context) string {
	actionIDPrefixes := shortenIDs(incompleteActionIDs(contexts...))
	outcomeIDPrefixes := shortenIDs(incompleteActionOutcomeIDs(contexts...))

	builder := new(strings.Builder)

	for _, context := range contexts {
		_, _ = fmt.Fprintf(builder, "## @%s:\n\n", strings.Title(context.Name))

		for _, action := range context.Actions {
			_, _ = fmt.Fprintf(builder,
				"- [ ] `0x%s` %s (`0x%s` %s)\n",
				actionIDPrefixes[action.ID], action.Description,
				outcomeIDPrefixes[action.OutcomeID], action.OutcomeTitle,
			)
		}

		builder.WriteString("\n\n")
	}

	return strings.TrimSpace(builder.String())
}

func incompleteActionOutcomeIDs(contexts ...*projections.Context) (ids_ []string) {
	unique := make(map[string]bool)
	for _, context := range contexts {
		for _, action := range context.Actions {
			unique[action.OutcomeID] = true
		}
	}
	for id := range unique {
		ids_ = append(ids_, id)
	}
	return ids_
}

func incompleteActionIDs(contexts ...*projections.Context) (ids_ []string) {
	unique := make(map[string]bool)
	for _, context := range contexts {
		for _, action := range context.Actions {
			unique[action.ID] = true
		}
	}
	for id := range unique {
		ids_ = append(ids_, id)
	}
	return ids_
}
