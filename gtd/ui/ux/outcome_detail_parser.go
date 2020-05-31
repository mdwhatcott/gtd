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

	handler   core.Handler
	scanner   *bufio.Scanner
	line      string
	words     []string
	parseLine func()

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
	if this.line == "# {TITLE}" {
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
		action := &commands.UpdateActionDescription{
			OutcomeID:          this.outcomeID,
			ActionID:           actionID,
			UpdatedDescription: this.parseActionDescription(),
		}
		this.handle(action)
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
func (this *OutcomeDetailParser) isExistingAction() bool {
	return this.actionLineHasID() && this.actionIDs[this.extractActionID()]
}
func (this *OutcomeDetailParser) actionLineHasID() bool {
	return strings.Contains(this.line, " `0x")
}
func (this *OutcomeDetailParser) extractActionID() string {
	start := strings.Index(this.line, " `0x") + 4
	end := start + strings.Index(this.line[start:], "` ")
	prefix := this.line[start:end]
	for id := range this.actionIDs {
		if strings.HasPrefix(id, prefix) {
			return id
		}
	}
	return prefix
}
func (this *OutcomeDetailParser) parseActionStatus(actionID string) {
	if this.isCompletedAction() {
		this.handle(&commands.MarkActionStatusComplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isIncompleteAction() {
		this.handle(&commands.MarkActionStatusIncomplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isLatentAction() {
		this.handle(&commands.MarkActionStatusLatent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	}
}
func (this *OutcomeDetailParser) parseActionStrategy(actionID string) {
	if this.isConcurrentAction() {
		this.handle(&commands.MarkActionStrategyConcurrent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isParallelAction() {
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
		open := strings.Index(this.line, " `0x") + 4
		end := open + strings.Index(this.line[open:], "` ")
		start = end + 1
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
