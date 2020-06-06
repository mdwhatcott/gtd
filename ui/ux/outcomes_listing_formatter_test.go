package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/core/projections"
)

func TestOutcomesListingFormatter(t *testing.T) {
	gunit.Run(new(OutcomesListingFormatter), t)
}

type OutcomesListingFormatter struct {
	*gunit.Fixture
}

func (this *OutcomesListingFormatter) Test() {
	listing := projections.OutcomesListing{
		Fixed: []*projections.OutcomesListingItem{
			{ID: "000111", Title: "1"},
			{ID: "000222", Title: "2"},
		},
		Deferred: []*projections.OutcomesListingItem{
			{ID: "000333", Title: "3"},
			{ID: "000444", Title: "4"},
		},
		Uncertain: []*projections.OutcomesListingItem{
			{ID: "000555", Title: "5"},
			{ID: "000666", Title: "6"},
		},
		Abandoned: []*projections.OutcomesListingItem{
			{ID: "000777", Title: "7"},
			{ID: "000888", Title: "8"},
		},
	}

	result := FormatOutcomesListing(listing)
	this.So(result, should.Equal, strings.Join([]string{
		"## Realized:",
		"",
		"",
		"## Fixed:",
		"",
		"- `0x0001` 1",
		"- `0x0002` 2",
		"",
		"",
		"## Deferred:",
		"",
		"- `0x0003` 3",
		"- `0x0004` 4",
		"",
		"",
		"## Uncertain:",
		"",
		"- `0x0005` 5",
		"- `0x0006` 6",
		"",
		"",
		"## Abandoned:",
		"",
		"- `0x0007` 7",
		"- `0x0008` 8",
		"",
		"",
		"## Deleted:",
		"",
	}, "\n"))
}
