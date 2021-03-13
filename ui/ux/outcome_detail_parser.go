package ux

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/mdwhatcott/gtd/v3/core"
	"github.com/mdwhatcott/gtd/v3/core/commands"
	"github.com/mdwhatcott/gtd/v3/core/projections"
)

type OutcomeDetailParser struct {
	handler    core.Handler
	projection projections.OutcomeDetails
	scanner    *bufio.Scanner
	line       string
	words      []string
	parseLine  func()

	outcomeID          string
	description        *strings.Builder
	actionIDs          map[string]bool
	updatedActionOrder []string
}

func NewOutcomeDetailParser(
	handler core.Handler,
	outcomeID string,
	projection projections.OutcomeDetails,
	modifiedContent string,
) *OutcomeDetailParser {
	UX := &OutcomeDetailParser{
		handler:     handler,
		projection:  projection,
		scanner:     bufio.NewScanner(strings.NewReader(modifiedContent)),
		outcomeID:   outcomeID,
		actionIDs:   actionIDMap(projection.Actions),
		description: new(strings.Builder),
	}
	UX.parseLine = UX.parseTitleLine
	return UX
}

func actionIDMap(actions []*projections.ActionDetails) map[string]bool {
	SET := make(map[string]bool)
	for _, action := range actions {
		SET[action.ID] = true
	}
	return SET
}

func (this *OutcomeDetailParser) handle(instructions ...interface{}) {
	handle(this.handler, instructions...)
}

func (this *OutcomeDetailParser) Parse() error {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		this.words = strings.Fields(this.line)
		this.parseLine()
		if this.outcomeID == "" {
			break
		}
	}
	this.reorderActions()
	this.deleteRemovedActions()
	this.updateDescription()
	return nil
}

func (this *OutcomeDetailParser) parseTitleLine() {
	if !strings.HasPrefix(this.line, "# ") {
		return
	}
	if this.outcomeID == "" {
		COMMAND := &commands.TrackOutcome{Title: this.line[2:]}
		this.handle(COMMAND)
		this.outcomeID = COMMAND.Result.ID
		this.handle(&commands.DeclareOutcomeFixed{OutcomeID: this.outcomeID})
	} else {
		TITLE := this.line[2:]
		if TITLE != this.projection.Title {
			this.handle(&commands.UpdateOutcomeTitle{OutcomeID: this.outcomeID, UpdatedTitle: TITLE})
		}
	}

	this.parseLine = this.parseExplanationLine
}
func (this *OutcomeDetailParser) parseExplanationLine() {
	if this.line == "## Actions:" {
		this.parseLine = this.parseActionLine
		return
	}
	if this.line == "" {
		return
	}
	if strings.HasPrefix(this.line, "> ") {
		explanation := this.line[2:]
		if explanation != this.projection.Explanation {
			this.handle(&commands.UpdateOutcomeExplanation{
				OutcomeID:          this.outcomeID,
				UpdatedExplanation: explanation,
			})
		}
	}
}
func (this *OutcomeDetailParser) parseActionLine() {
	if this.line == "## Support Materials:" {
		this.parseLine = this.parseOutcomeDescriptionLine
	}
	if this.line == "" {
		return
	}
	if !this.currentLineIsAction() {
		return
	}

	this.parseAction()
}
func (this *OutcomeDetailParser) parseAction() {
	var ACTION_ID string
	if this.isExistingAction() {
		ACTION_ID = this.extractActionID()
		DESCRIPTION := this.parseActionDescription()
		if this.descriptionModified(ACTION_ID, DESCRIPTION) {
			this.handle(&commands.UpdateActionDescription{
				OutcomeID:          this.outcomeID,
				ActionID:           ACTION_ID,
				UpdatedDescription: DESCRIPTION,
			})
		}
		this.actionIDs[ACTION_ID] = false
	} else {
		COMMAND := &commands.TrackAction{
			OutcomeID:   this.outcomeID,
			Description: this.parseActionDescription(),
		}
		this.handle(COMMAND)
		ACTION_ID = COMMAND.Result.ID
	}

	this.updatedActionOrder = append(this.updatedActionOrder, ACTION_ID)
	this.parseActionStrategy(ACTION_ID)
	this.parseActionStatus(ACTION_ID)
}
func (this *OutcomeDetailParser) descriptionModified(id string, updatedDescription string) bool {
	for _, ACTION := range this.projection.Actions {
		if ACTION.ID == id {
			return updatedDescription != ACTION.Description
		}
	}
	return false
}
func (this *OutcomeDetailParser) isExistingAction() bool {
	return this.actionLineHasID() && this.actionIDs[this.extractActionID()]
}
func (this *OutcomeDetailParser) actionLineHasID() bool {
	return strings.Contains(this.line, " `0x")
}
func (this *OutcomeDetailParser) extractActionID() string {
	PREFIX := between(this.line, " `0x", "` ")
	for ID := range this.actionIDs {
		if strings.HasPrefix(ID, PREFIX) {
			return ID
		}
	}
	return PREFIX
}
func (this *OutcomeDetailParser) parseActionStatus(actionID string) {
	if this.isCompletedAction() && !this.actionStatusPreviouslyMatched(actionID, core.ActionStatusComplete) {
		this.handle(&commands.MarkActionStatusComplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isIncompleteAction() && !this.actionStatusPreviouslyMatched(actionID, core.ActionStatusIncomplete) {
		this.handle(&commands.MarkActionStatusIncomplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isLatentAction() && !this.actionStatusPreviouslyMatched(actionID, core.ActionStatusLatent) {
		this.handle(&commands.MarkActionStatusLatent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	}
}
func (this *OutcomeDetailParser) actionStatusPreviouslyMatched(actionID string, previousStatus core.ActionStatus) bool {
	for _, ACTION := range this.projection.Actions {
		if ACTION.ID == actionID {
			return ACTION.Status == previousStatus
		}
	}
	return false

}
func (this *OutcomeDetailParser) parseActionStrategy(actionID string) {
	if this.isConcurrentAction() && !this.actionStrategyPreviouslyMatched(actionID, core.ActionStrategyConcurrent) {
		this.handle(&commands.MarkActionStrategyConcurrent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isParallelAction() && !this.actionStrategyPreviouslyMatched(actionID, core.ActionStrategySequential) {
		this.handle(&commands.MarkActionStrategySequential{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	}
}
func (this *OutcomeDetailParser) currentLineIsAction() bool {
	return this.isConcurrentAction() || this.isParallelAction()
}
func (this *OutcomeDetailParser) isConcurrentAction() bool {
	return this.words[0] == "-"
}
func (this *OutcomeDetailParser) actionStrategyPreviouslyMatched(actionID string, previousStrategy core.ActionStrategy) bool {
	for _, ACTION := range this.projection.Actions {
		if ACTION.ID == actionID {
			return ACTION.Strategy == previousStrategy
		}
	}
	return false
}
func (this *OutcomeDetailParser) isParallelAction() bool {
	FIRST := this.words[0]
	FIRST = strings.TrimRight(FIRST, ".")
	number, _ := strconv.Atoi(FIRST)
	return number > 0
}
func (this *OutcomeDetailParser) isCompletedAction() bool {
	return this.words[1] == "[X]" || this.words[1] == "[x]"
}
func (this *OutcomeDetailParser) isLatentAction() bool {
	return this.words[1] == "[?]"
}
func (this *OutcomeDetailParser) isIncompleteAction() bool {
	return this.words[1] == "[]" || (this.words[1] == "[" && this.words[2] == "]")
}
func (this *OutcomeDetailParser) parseActionDescription() string {
	START := strings.Index(this.line, "]") + 1
	if this.actionLineHasID() {
		START = strings.Index(this.line, "` ") + 1
	}
	return strings.TrimSpace(this.line[START:])
}
func (this *OutcomeDetailParser) parseOutcomeDescriptionLine() {
	_, _ = io.WriteString(this.description, this.line)
	_, _ = io.WriteString(this.description, "\n")
}
func (this *OutcomeDetailParser) updateDescription() {
	DESCRIPTION := strings.TrimSpace(this.description.String())
	if DESCRIPTION != this.projection.Description {
		this.handle(&commands.UpdateOutcomeDescription{
			OutcomeID:          this.outcomeID,
			UpdatedDescription: DESCRIPTION,
		})
	}
}

func (this *OutcomeDetailParser) deleteRemovedActions() {
	for ID, REMOVED := range this.actionIDs {
		if REMOVED {
			this.handle(&commands.DeleteAction{
				OutcomeID: this.outcomeID,
				ActionID:  ID,
			})
		}
	}
}

func (this *OutcomeDetailParser) reorderActions() {
	if len(this.updatedActionOrder) == 0 {
		return
	}
	if len(this.projection.Actions) == 0 {
		return
	}
	var CURRENT []string
	for _, action := range this.projection.Actions {
		CURRENT = append(CURRENT, action.ID)
	}
	A := strings.Join(CURRENT, "|")
	B := strings.Join(this.updatedActionOrder, "|")
	if A == B {
		return
	}
	if strings.HasPrefix(B, A) { // only post-inserts happened, no need to reorder
		return
	}

	this.handle(&commands.ReorderActions{
		OutcomeID:    this.outcomeID,
		ReorderedIDs: this.updatedActionOrder,
	})
}
