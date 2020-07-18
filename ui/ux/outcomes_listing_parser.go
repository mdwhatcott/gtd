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

func (this *OutcomesListingParser) Parse() (edits_ []string) {
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
	STATUS, FOUND := headerToStatus[this.line]
	if FOUND {
		this.status = STATUS
	}
	return FOUND
}

func (this *OutcomesListingParser) parseLine() {
	ID := between(this.line, "`0x", "`")
	if ID == "" {
		this.trackOutcome()
	} else {
		this.updateOutcome(ID)
	}
}

func (this *OutcomesListingParser) trackOutcome() {
	COMMAND := &commands.TrackOutcome{Title: strings.TrimSpace(strings.TrimLeft(this.line, "-"))}
	this.handle(COMMAND)
	this.changeStatus(COMMAND.Result.ID)
	this.edits = append(this.edits, COMMAND.Result.ID)
}

func (this *OutcomesListingParser) updateOutcome(idPrefix string) {
	ID, STATUS := this.findID(idPrefix)
	if STATUS != this.status {
		this.changeStatus(ID)
	}
	if strings.HasPrefix(this.line, "\t") {
		this.edits = append(this.edits, ID)
	}
}

func (this *OutcomesListingParser) findID(prefix string) (string, core.OutcomeStatus) {
	for _, OUTCOME := range this.combinedListings() {
		if strings.HasPrefix(OUTCOME.ID, prefix) {
			return OUTCOME.ID, OUTCOME.Status
		}
	}
	return prefix, ""
}

func (this *OutcomesListingParser) combinedListings() (all_ []*projections.OutcomesListingItem) {
	all_ = append(all_, this.listings.Fixed...)
	all_ = append(all_, this.listings.Deferred...)
	all_ = append(all_, this.listings.Uncertain...)
	all_ = append(all_, this.listings.Abandoned...)
	all_ = append(all_, this.listings.Realized...)
	return all_
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
