package ux

import (
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/projections"
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
	this.listings.Fixed = append(this.listings.Fixed, &projections.OutcomesListingItem{ID: "000111", Title: "1"})
	this.listings.Deferred = append(this.listings.Deferred, &projections.OutcomesListingItem{ID: "000222", Title: "2"})
	this.listings.Uncertain = append(this.listings.Uncertain, &projections.OutcomesListingItem{ID: "000333", Title: "3"})
	this.listings.Abandoned = append(this.listings.Abandoned,
		&projections.OutcomesListingItem{ID: "000444", Title: "4"},
		&projections.OutcomesListingItem{ID: "000555", Title: "5"},
		&projections.OutcomesListingItem{ID: "000666", Title: "6"},
	)
	this.content = strings.Join([]string{
		"## Realized:",
		"- `0x0001` 1",

		"## Fixed:",
		"- `0x0002` 2",

		"## Deferred:",
		"- `0x0003` 3",

		"## Uncertain:",
		"- `0x0004` 4",

		"## Abandoned:",
		"- `0x0005` 5",

		"## Deleted:",
		"- `0x0006` 6",
	}, "\n")

	parser := NewOutcomesListingParser(this.handler, this.listings, this.content)

	err := parser.Parse()

	this.So(err, should.BeNil)
	this.So(this.handler.handled, should.Resemble, []interface{}{
		&commands.DeclareOutcomeRealized{OutcomeID: "000111"},
		&commands.DeclareOutcomeFixed{OutcomeID: "000222"},
		&commands.DeclareOutcomeDeferred{OutcomeID: "000333"},
		&commands.DeclareOutcomeUncertain{OutcomeID: "000444"},
		&commands.DeclareOutcomeAbandoned{OutcomeID: "000555"},
		&commands.DeleteOutcome{OutcomeID: "000666"},
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
}
