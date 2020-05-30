package ux

import (
	"bufio"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

type OutcomesListingParser struct {
	handler  core.Handler
	listings projections.OutcomesListing
	scanner  *bufio.Scanner
	status   core.OutcomeStatus
	line     string
	words    []string
}

func NewOutcomesListingParser(
	handler core.Handler,
	listings projections.OutcomesListing,
	content string,
) *OutcomesListingParser {
	return &OutcomesListingParser{
		handler:  handler,
		listings: listings,
		scanner:  bufio.NewScanner(strings.NewReader(content)),
	}
}

func (this *OutcomesListingParser) Parse() error {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		if len(this.line) == 0 {
			continue
		}

		this.words = strings.Fields(this.line)
		if !this.updateStatus() {
			this.parseLine()
		}
	}
	return nil
}

func (this *OutcomesListingParser) updateStatus() bool {
	status, found := headerToStatus[this.line]
	if found {
		this.status = status
	}
	return found
}

func (this *OutcomesListingParser) parseLine() {
	idPrefix := this.words[1]
	idPrefix = strings.TrimPrefix(idPrefix, "`0x")
	idPrefix = strings.TrimSuffix(idPrefix, "`")
	this.sendInstruction(this.findID(idPrefix))
}

func (this *OutcomesListingParser) findID(prefix string) string {
	var all []*projections.OutcomesListingItem
	all = append(all, this.listings.Fixed...)
	all = append(all, this.listings.Deferred...)
	all = append(all, this.listings.Uncertain...)
	all = append(all, this.listings.Abandoned...)
	for _, outcome := range all {
		if strings.HasPrefix(outcome.ID, prefix) {
			return outcome.ID
		}
	}
	return prefix
}

func (this *OutcomesListingParser) sendInstruction(id string) {
	switch this.status {

	case core.OutcomeStatusRealized:
		this.handler.Handle(&commands.DeclareOutcomeRealized{OutcomeID: id})

	case core.OutcomeStatusFixed:
		this.handler.Handle(&commands.DeclareOutcomeFixed{OutcomeID: id})

	case core.OutcomeStatusDeferred:
		this.handler.Handle(&commands.DeclareOutcomeDeferred{OutcomeID: id})

	case core.OutcomeStatusUncertain:
		this.handler.Handle(&commands.DeclareOutcomeUncertain{OutcomeID: id})

	case core.OutcomeStatusAbandoned:
		this.handler.Handle(&commands.DeclareOutcomeAbandoned{OutcomeID: id})

	case "":
		this.handler.Handle(&commands.DeleteOutcome{OutcomeID: id})
	}
}

var headerToStatus = map[string]core.OutcomeStatus{
	"## Realized:":  core.OutcomeStatusRealized,
	"## Fixed:":     core.OutcomeStatusFixed,
	"## Deferred:":  core.OutcomeStatusDeferred,
	"## Uncertain:": core.OutcomeStatusUncertain,
	"## Abandoned:": core.OutcomeStatusAbandoned,
	"## Deleted:":   "",
}
