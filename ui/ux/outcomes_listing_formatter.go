package ux

import (
	"strings"

	"github.com/mdwhatcott/gtd/core/projections"
)

func FormatOutcomesListing(listing projections.OutcomesListing) string {
	builder := new(strings.Builder)
	builder.WriteString(composeListing("## Fixed:", listing.Fixed))
	builder.WriteString(composeListing("## Deferred:", listing.Deferred))
	builder.WriteString(composeListing("## Uncertain:", listing.Uncertain))
	builder.WriteString(composeListing("## Abandoned:", listing.Abandoned))
	builder.WriteString(composeListing("## Realized:", listing.Realized))
	builder.WriteString(composeListing("## Deleted:", nil))
	return strings.TrimSpace(builder.String()) + "\n"
}

func composeListing(header string, items []*projections.OutcomesListingItem) string {
	builder := new(strings.Builder)
	builder.WriteString(header)
	builder.WriteString("\n")
	builder.WriteString("\n")
	IDs := shortenIDs(outcomesListingIDs(items))
	for _, outcome := range items {
		builder.WriteString("- `0x")
		builder.WriteString(IDs[outcome.ID])
		builder.WriteString("` ")
		builder.WriteString(outcome.Title)
		builder.WriteString("\n")
	}
	if len(items) > 0 {
		builder.WriteString("\n")
	}
	builder.WriteString("\n")
	return builder.String()
}

func outcomesListingIDs(listing []*projections.OutcomesListingItem) (ids_ []string) {
	for _, item := range listing {
		ids_ = append(ids_, item.ID)
	}
	return ids_
}

func shortenIDs(_ids []string) (fullToPrefix map[string]string) {
	if len(_ids) == 0 {
		return nil
	}
	for length := minIDPrefixLength; length < len(_ids[0]); length++ {
		actionIDs := make(map[string]bool)
		fullToPrefix = make(map[string]string)
		for _, ID := range _ids {
			prefix := ID[:length]
			actionIDs[prefix] = true
			fullToPrefix[ID] = prefix
		}
		if len(actionIDs) == len(_ids) {
			break
		}
	}
	return fullToPrefix
}

const minIDPrefixLength = 4
