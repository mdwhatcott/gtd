package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

func TestOutcomesListingParserFixture(t *testing.T) {
	gunit.Run(new(OutcomesListingParserFixture), t)
}

type OutcomesListingParserFixture struct {
	*gunit.Fixture
	handler  *FakeOutcomesListingParserFakeHandler
	listings projections.OutcomesListing
	content  string
}

func (this *OutcomesListingParserFixture) Setup() {
	this.handler = NewFakeOutcomesListingParserFakeHandler()
}

func (this *OutcomesListingParserFixture) Test() {
	this.listings.Fixed = append(this.listings.Fixed,
		&projections.OutcomesListingItem{ID: "000111", Title: "1", Status: core.OutcomeStatusFixed},
		&projections.OutcomesListingItem{ID: "000222", Title: "2", Status: core.OutcomeStatusFixed},
		&projections.OutcomesListingItem{ID: "000333", Title: "3", Status: core.OutcomeStatusFixed},
	)
	this.listings.Deferred = append(this.listings.Deferred,
		&projections.OutcomesListingItem{ID: "000444", Title: "4", Status: core.OutcomeStatusDeferred},
	)
	this.listings.Uncertain = append(this.listings.Uncertain,
		&projections.OutcomesListingItem{ID: "000555", Title: "5", Status: core.OutcomeStatusUncertain},
	)
	this.listings.Abandoned = append(this.listings.Abandoned,
		&projections.OutcomesListingItem{ID: "000666", Title: "6", Status: core.OutcomeStatusAbandoned},
		&projections.OutcomesListingItem{ID: "000777", Title: "7", Status: core.OutcomeStatusAbandoned},
		&projections.OutcomesListingItem{ID: "000888", Title: "8", Status: core.OutcomeStatusAbandoned},
	)
	this.listings.Realized = append(this.listings.Realized,
		&projections.OutcomesListingItem{ID: "000999", Title: "9", Status: core.OutcomeStatusRealized},
	)
	this.content = strings.Join([]string{
		"## Fixed:",
		"\t- `0x0002` 2",
		"- `0x0003` 3",

		"## Deferred:",
		"- `0x0005` 5",

		"## Uncertain:",
		"- `0x0004` 4",
		"- All New Outcome!",

		"## Abandoned:",
		"- `0x0006` 6",

		"## Realized:",
		"- `0x0001` 1",

		"## Deleted:",
		"- `0x0007` 7",
		"- `0x0008` 8",
	}, "\n")

	parser := NewOutcomesListingParser(this.handler, this.listings, this.content)

	requestedEdits := parser.Parse()

	this.So(requestedEdits, should.Resemble, []string{"000222", "0042"})
	this.So(this.handler.handled, should.Resemble, []interface{}{
		//&commands.DeclareOutcomeFixed{OutcomeID: "000222"}, // unchanged
		//&commands.DeclareOutcomeFixed{OutcomeID: "000333"}, // unchanged
		&commands.DeclareOutcomeDeferred{OutcomeID: "000555"},
		&commands.DeclareOutcomeUncertain{OutcomeID: "000444"},
		&commands.TrackOutcome{Title: "All New Outcome!", Result: commands.CreateResult{ID: "0042"}},
		&commands.DeclareOutcomeUncertain{OutcomeID: "0042"},
		//&commands.DeclareOutcomeAbandoned{OutcomeID: "000666"}, // unchanged
		&commands.DeclareOutcomeRealized{OutcomeID: "000111"},
		&commands.DeleteOutcome{OutcomeID: "000777"},
		&commands.DeleteOutcome{OutcomeID: "000888"},
	})
}

//////////////////////////////////////////////////////////////////////

type FakeOutcomesListingParserFakeHandler struct {
	handled []interface{}
}

func NewFakeOutcomesListingParserFakeHandler() *FakeOutcomesListingParserFakeHandler {
	return &FakeOutcomesListingParserFakeHandler{}
}

func (this *FakeOutcomesListingParserFakeHandler) Handle(messages ...interface{}) {
	this.handled = append(this.handled, messages...)
	for _, message := range messages {
		switch message := message.(type) {
		case *commands.TrackOutcome:
			message.Result.ID = "0042"
		}
	}
}
