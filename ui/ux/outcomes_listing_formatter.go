package ux

import (
	"strings"

	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func FormatOutcomesListing(listing projections.OutcomesListing) string {
	BUILDER := new(strings.Builder)
	BUILDER.WriteString(composeListing("## Fixed:", listing.Fixed))
	BUILDER.WriteString(composeListing("## Deferred:", listing.Deferred))
	BUILDER.WriteString(composeListing("## Uncertain:", listing.Uncertain))
	BUILDER.WriteString(composeListing("## Abandoned:", listing.Abandoned))
	BUILDER.WriteString(composeListing("## Realized:", listing.Realized))
	BUILDER.WriteString(composeListing("## Deleted:", nil))
	return strings.TrimSpace(BUILDER.String()) + "\n"
}

func composeListing(header string, items []*projections.OutcomesListingItem) string {
	BUILDER := new(strings.Builder)
	BUILDER.WriteString(header)
	BUILDER.WriteString("\n")
	BUILDER.WriteString("\n")
	IDs := shortenIDs(outcomesListingIDs(items))
	for _, OUTCOME := range items {
		BUILDER.WriteString("- `0x")
		BUILDER.WriteString(IDs[OUTCOME.ID])
		BUILDER.WriteString("` ")
		BUILDER.WriteString(OUTCOME.Title)
		BUILDER.WriteString("\n")
	}
	if len(items) > 0 {
		BUILDER.WriteString("\n")
	}
	BUILDER.WriteString("\n")
	return BUILDER.String()
}

func outcomesListingIDs(listing []*projections.OutcomesListingItem) (ids_ []string) {
	for _, ITEM := range listing {
		ids_ = append(ids_, ITEM.ID)
	}
	return ids_
}

func shortenIDs(ids []string) (fullToPrefix_ map[string]string) {
	if len(ids) == 0 {
		return nil
	}
	for LENGTH := minIDPrefixLength; LENGTH < len(ids[0]); LENGTH++ {
		ACTIONS := make(map[string]bool)
		fullToPrefix_ = make(map[string]string)
		for _, ID := range ids {
			PREFIX := ID[:LENGTH]
			ACTIONS[PREFIX] = true
			fullToPrefix_[ID] = PREFIX
		}
		if len(ACTIONS) == len(ids) {
			break
		}
	}
	return fullToPrefix_
}

const minIDPrefixLength = 4
