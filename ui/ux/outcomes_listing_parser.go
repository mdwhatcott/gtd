package ux

import (
	"bufio"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

type OutcomesListingParser struct {
	handler  core.Handler
	listings projections.OutcomesListing
	scanner  *bufio.Scanner
	status   core.OutcomeStatus
	line     string
	edits    []string
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

func (this *OutcomesListingParser) handle(instruction interface{}) {
	handle(this.handler, instruction)
}

func (this *OutcomesListingParser) Parse() (edits []string) {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		if len(this.line) == 0 {
			continue
		}

		if !this.updateStatus() {
			this.parseLine()
		}
	}
	return this.edits
}

func (this *OutcomesListingParser) updateStatus() bool {
	status, found := headerToStatus[this.line]
	if found {
		this.status = status
	}
	return found
}

func (this *OutcomesListingParser) parseLine() {
	idPrefix := between(this.line, "`0x", "`")
	if idPrefix == "" {
		this.trackOutcome()
	} else {
		this.updateOutcome(idPrefix)
	}
}

func (this *OutcomesListingParser) trackOutcome() {
	command := &commands.TrackOutcome{Title: strings.TrimSpace(strings.TrimLeft(this.line, "-"))}
	this.handle(command)
	this.changeStatus(command.Result.ID)
	this.edits = append(this.edits, command.Result.ID)
}

func (this *OutcomesListingParser) updateOutcome(idPrefix string) {
	id, status := this.findID(idPrefix)
	if status != this.status {
		this.changeStatus(id)
	}
	if strings.HasPrefix(this.line, "\t") {
		this.edits = append(this.edits, id)
	}
}

func (this *OutcomesListingParser) findID(prefix string) (string, core.OutcomeStatus) {
	for _, outcome := range this.combinedListings() {
		if strings.HasPrefix(outcome.ID, prefix) {
			return outcome.ID, outcome.Status
		}
	}
	return prefix, ""
}

func (this *OutcomesListingParser) combinedListings() (all []*projections.OutcomesListingItem) {
	all = append(all, this.listings.Fixed...)
	all = append(all, this.listings.Deferred...)
	all = append(all, this.listings.Uncertain...)
	all = append(all, this.listings.Abandoned...)
	all = append(all, this.listings.Realized...)
	return all
}

func (this *OutcomesListingParser) changeStatus(id string) {
	switch this.status {

	case core.OutcomeStatusRealized:
		this.handle(&commands.DeclareOutcomeRealized{OutcomeID: id})

	case core.OutcomeStatusFixed:
		this.handle(&commands.DeclareOutcomeFixed{OutcomeID: id})

	case core.OutcomeStatusDeferred:
		this.handle(&commands.DeclareOutcomeDeferred{OutcomeID: id})

	case core.OutcomeStatusUncertain:
		this.handle(&commands.DeclareOutcomeUncertain{OutcomeID: id})

	case core.OutcomeStatusAbandoned:
		this.handle(&commands.DeclareOutcomeAbandoned{OutcomeID: id})

	case "":
		this.handle(&commands.DeleteOutcome{OutcomeID: id})
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
