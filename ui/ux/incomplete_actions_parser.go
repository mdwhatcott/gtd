package ux

import (
	"bufio"
	"strings"

	"github.com/mdwhatcott/gtd/core"
	"github.com/mdwhatcott/gtd/core/commands"
	"github.com/mdwhatcott/gtd/core/projections"
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
	_, outcomeID := this.identifyIDs()
	if outcomeID == "" {
		return
	}
	this.edits = append(this.edits, outcomeID)
}

func (this *IncompleteActionsParser) markActionLatent() {
	actionID, outcomeID := this.identifyIDs()
	this.handler.Handle(&commands.MarkActionStatusLatent{
		OutcomeID: outcomeID,
		ActionID:  actionID,
	})
}

func (this *IncompleteActionsParser) markActionComplete() {
	actionID, outcomeID := this.identifyIDs()
	this.handler.Handle(&commands.MarkActionStatusComplete{
		OutcomeID: outcomeID,
		ActionID:  actionID,
	})
}

func (this *IncompleteActionsParser) identifyIDs() (actionID, outcomeID string) {
	actionID = this.findID(between(this.line, "`0x", "`"))
	outcomeID = this.findID(between(this.line[strings.Index(this.line, "("):], "(`0x", "`"))
	return actionID, outcomeID
}

func (this *IncompleteActionsParser) findID(prefix string) string {
	for _, context := range this.contexts {
		for _, action := range context.Actions {
			if strings.HasPrefix(action.ID, prefix) {
				return action.ID
			}
			if strings.HasPrefix(action.OutcomeID, prefix) {
				return action.OutcomeID
			}
		}
	}
	return ""
}

func between(value, left, right string) string {
	index1 := strings.Index(value, left)
	if index1 < 0 {
		return ""
	}
	start := index1 + len(left)

	index2 := strings.Index(value[start:], right)
	if index2 < 0 {
		return ""
	}
	stop := start + index2

	return value[start:stop]
}
