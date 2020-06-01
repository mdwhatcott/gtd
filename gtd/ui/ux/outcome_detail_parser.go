package ux

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/smartystreets/logging"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/core/projections"
)

type OutcomeDetailParser struct {
	log *logging.Logger

	handler    core.Handler
	projection projections.OutcomeDetails
	scanner    *bufio.Scanner
	line       string
	words      []string
	parseLine  func()

	outcomeID   string
	description *strings.Builder
	actionIDs   map[string]bool
}

func NewOutcomeDetailParser(
	handler core.Handler,
	outcomeID string,
	projection projections.OutcomeDetails,
	modifiedContent string,
) *OutcomeDetailParser {
	ux := &OutcomeDetailParser{
		handler:     handler,
		projection:  projection,
		scanner:     bufio.NewScanner(strings.NewReader(modifiedContent)),
		outcomeID:   outcomeID,
		actionIDs:   actionIDMap(projection.Actions),
		description: new(strings.Builder),
	}
	ux.parseLine = ux.parseTitleLine
	return ux
}

func actionIDMap(actions []*projections.ActionDetails) map[string]bool {
	set := make(map[string]bool)
	for _, action := range actions {
		set[action.ID] = true
	}
	return set
}

func (this *OutcomeDetailParser) handle(instructions ...interface{}) {
	this.handler.Handle(instructions...)
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
	this.deleteRemovedActions()
	this.updateDescription()
	return nil
}

func (this *OutcomeDetailParser) parseTitleLine() {
	if !strings.HasPrefix(this.line, "# ") {
		return
	}
	if this.outcomeID == "" {
		command := &commands.TrackOutcome{Title: this.line[2:]}
		this.handle(command)
		this.outcomeID = command.Result.ID
		this.handle(&commands.DeclareOutcomeFixed{OutcomeID: this.outcomeID})
	} else {
		this.handle(&commands.UpdateOutcomeTitle{UpdatedTitle: this.line[2:]})
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
		this.handle(&commands.UpdateOutcomeExplanation{
			OutcomeID:          this.outcomeID,
			UpdatedExplanation: this.line[2:],
		})
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
	var actionID string
	if this.isExistingAction() {
		actionID = this.extractActionID()
		description := this.parseActionDescription()
		if this.descriptionModified(actionID, description) {
			this.handle(&commands.UpdateActionDescription{
				OutcomeID:          this.outcomeID,
				ActionID:           actionID,
				UpdatedDescription: description,
			})
		}
		this.actionIDs[actionID] = false
	} else {
		action := &commands.TrackAction{
			OutcomeID:   this.outcomeID,
			Description: this.parseActionDescription(),
		}
		this.handle(action)
		actionID = action.Result.ID
	}

	this.parseActionStrategy(actionID)
	this.parseActionStatus(actionID)
}
func (this *OutcomeDetailParser) descriptionModified(id string, updatedDescription string) bool {
	for _, action := range this.projection.Actions {
		if action.ID == id {
			return updatedDescription != action.Description
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
	prefix := between(this.line, " `0x", "` ")
	for id := range this.actionIDs {
		if strings.HasPrefix(id, prefix) {
			return id
		}
	}
	return prefix
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
	for _, action := range this.projection.Actions {
		if action.ID == actionID {
			return action.Status == previousStatus
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
	for _, action := range this.projection.Actions {
		if action.ID == actionID {
			return action.Strategy == previousStrategy
		}
	}
	return false
}
func (this *OutcomeDetailParser) isParallelAction() bool {
	first := this.words[0]
	first = strings.TrimRight(first, ".")
	number, _ := strconv.Atoi(first)
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
	start := strings.Index(this.line, "]") + 1
	if this.actionLineHasID() {
		start = strings.Index(this.line, "` ") + 1
	}
	return strings.TrimSpace(this.line[start:])
}
func (this *OutcomeDetailParser) parseOutcomeDescriptionLine() {
	_, _ = io.WriteString(this.description, this.line)
	_, _ = io.WriteString(this.description, "\n")
}
func (this *OutcomeDetailParser) updateDescription() {
	if this.description.Len() > 0 {
		this.handle(&commands.UpdateOutcomeDescription{
			OutcomeID:          this.outcomeID,
			UpdatedDescription: strings.TrimSpace(this.description.String()),
		})
	}
}

func (this *OutcomeDetailParser) deleteRemovedActions() {
	for id, removed := range this.actionIDs {
		if removed {
			this.handle(&commands.DeleteAction{
				OutcomeID: this.outcomeID,
				ActionID:  id,
			})
		}
	}
}
