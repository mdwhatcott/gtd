package ux

import (
	"bufio"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

type IncompleteActionsParser struct {
	handler  core.Handler
	contexts []*projections.Context
	scanner  *bufio.Scanner
	line     string
	edits    []string
}

func NewIncompleteActionsParser(handler core.Handler, content string, contexts ...*projections.Context) *IncompleteActionsParser {
	return &IncompleteActionsParser{
		handler:  handler,
		contexts: contexts,
		scanner:  bufio.NewScanner(strings.NewReader(content)),
	}
}

func (this *IncompleteActionsParser) handle(instructions ...interface{}) {
	handle(this.handler, instructions...)
}

func (this *IncompleteActionsParser) Parse() (edits []string) {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		if this.line == "" {
			continue
		}
		if this.lineIsContextHeader() {
			continue
		}

		if this.outcomeReviewRequested() {
			this.parseEdits()
		}
		if this.actionCompleted() {
			this.markActionComplete()
		} else if this.actionLatent() {
			this.markActionLatent()
		}
	}
	return this.edits
}

func (this *IncompleteActionsParser) lineIsContextHeader() bool {
	return strings.HasPrefix(this.line, "## ")
}

func (this *IncompleteActionsParser) outcomeReviewRequested() bool {
	return strings.HasPrefix(this.line, "\t- ")
}

func (this *IncompleteActionsParser) actionLatent() bool {
	return strings.HasPrefix(this.line, "- ") && strings.ToUpper(between(this.line, "[", "]")) == "?"
}

func (this *IncompleteActionsParser) actionCompleted() bool {
	return strings.HasPrefix(this.line, "- ") && strings.ToUpper(between(this.line, "[", "]")) == "X"
}

func (this *IncompleteActionsParser) parseEdits() {
	this.line = strings.TrimSpace(this.line)
	_, OUTCOME := this.identifyIDs()
	if OUTCOME == "" {
		return
	}
	this.edits = append(this.edits, OUTCOME)
}

func (this *IncompleteActionsParser) markActionLatent() {
	ACTION, OUTCOME := this.identifyIDs()
	this.handle(&commands.MarkActionStatusLatent{
		OutcomeID: OUTCOME,
		ActionID:  ACTION,
	})
}

func (this *IncompleteActionsParser) markActionComplete() {
	ACTION, OUTCOME := this.identifyIDs()
	this.handle(&commands.MarkActionStatusComplete{
		OutcomeID: OUTCOME,
		ActionID:  ACTION,
	})
}

func (this *IncompleteActionsParser) identifyIDs() (actionID_, outcomeID_ string) {
	actionID_ = this.findID(between(this.line, "`0x", "`"))
	outcomeID_ = this.findID(between(this.line[strings.Index(this.line, "("):], "(`0x", "`"))
	return actionID_, outcomeID_
}

func (this *IncompleteActionsParser) findID(prefix string) string {
	for _, CONTEXT := range this.contexts {
		for _, ACTION := range CONTEXT.Actions {
			if strings.HasPrefix(ACTION.ID, prefix) {
				return ACTION.ID
			}
			if strings.HasPrefix(ACTION.OutcomeID, prefix) {
				return ACTION.OutcomeID
			}
		}
	}
	return ""
}

func between(value, left, right string) string {
	INDEX1 := strings.Index(value, left)
	if INDEX1 < 0 {
		return ""
	}
	START := INDEX1 + len(left)

	INDEX2 := strings.Index(value[START:], right)
	if INDEX2 < 0 {
		return ""
	}
	STOP := START + INDEX2

	return value[START:STOP]
}
