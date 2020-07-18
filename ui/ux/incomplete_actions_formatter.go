package ux

import (
	"fmt"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func FormatIncompleteActions(contexts ...*projections.Context) string {
	ACTIONS := shortenIDs(map2Slice(uniqueActionIDs(contexts...)))
	OUTCOMES := shortenIDs(map2Slice(uniqueActionOutcomeIDs(contexts...)))
	BUILDER := new(strings.Builder)

	for _, CONTEXT := range contexts {
		_, _ = fmt.Fprintf(BUILDER, "## @%s:\n\n", strings.Title(CONTEXT.Name))

		for _, ACTION := range CONTEXT.Actions {
			_, _ = fmt.Fprintf(BUILDER,
				"- [ ] `0x%s` %s (`0x%s` %s)\n",
				ACTIONS[ACTION.ID], ACTION.Description,
				OUTCOMES[ACTION.OutcomeID], ACTION.OutcomeTitle,
			)
		}

		BUILDER.WriteString("\n\n")
	}

	return strings.TrimSpace(BUILDER.String())
}

func map2Slice(keys map[string]bool) (slice_ []string) {
	for ID := range keys {
		slice_ = append(slice_, ID)
	}
	return slice_
}
func uniqueActionOutcomeIDs(contexts ...*projections.Context) (unique_ map[string]bool) {
	unique_ = make(map[string]bool)
	for _, CONTEXT := range contexts {
		for _, ACTION := range CONTEXT.Actions {
			unique_[ACTION.OutcomeID] = true
		}
	}
	return unique_
}
func uniqueActionIDs(contexts ...*projections.Context) (unique_ map[string]bool) {
	unique_ = make(map[string]bool)
	for _, CONTEXT := range contexts {
		for _, ACTION := range CONTEXT.Actions {
			unique_[ACTION.ID] = true
		}
	}
	return unique_
}
