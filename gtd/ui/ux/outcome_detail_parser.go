package ux

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
)

type OutcomeDetailParser struct {
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
	actionIDs map[string]bool,
	modifiedContent string,
) *OutcomeDetailParser {
	ux := &OutcomeDetailParser{
		handler:     handler,
		scanner:     bufio.NewScanner(strings.NewReader(modifiedContent)),
		outcomeID:   outcomeID,
		actionIDs:   actionIDs,
		description: new(strings.Builder),
	}
	ux.parseLine = ux.parseTitleLine
	return ux
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
	if this.description.Len() > 0 {
		this.handler.Handle(&commands.UpdateOutcomeDescription{
			OutcomeID:          this.outcomeID,
			UpdatedDescription: this.description.String(),
		})
	}
	return nil
}

func (this *OutcomeDetailParser) parseTitleLine() {
	if !strings.HasPrefix(this.line, "# ") {
		return
	}
	if this.line == "# {TITLE}" {
		return
	}
	command := &commands.TrackOutcome{Title: this.line[2:]}
	this.handler.Handle(command)
	this.outcomeID = command.Result.ID
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
		this.handler.Handle(&commands.UpdateOutcomeExplanation{
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
	action := &commands.TrackAction{
		OutcomeID:   this.outcomeID,
		Description: this.parseActionDescription(),
	}
	this.handler.Handle(action)

	this.parseActionStrategy(action.Result.ID)
	this.parseActionStatus(action.Result.ID)
}
func (this *OutcomeDetailParser) parseActionStatus(actionID string) {
	if this.isCompletedTask() {
		this.handler.Handle(&commands.MarkActionStatusComplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isIncompleteTask() {
		this.handler.Handle(&commands.MarkActionStatusIncomplete{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isLatentTask() {
		this.handler.Handle(&commands.MarkActionStatusLatent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	}
}
func (this *OutcomeDetailParser) parseActionStrategy(actionID string) {
	if this.isConcurrentTask() {
		this.handler.Handle(&commands.MarkActionStrategyConcurrent{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	} else if this.isParallelTask() {
		this.handler.Handle(&commands.MarkActionStrategySequential{
			OutcomeID: this.outcomeID,
			ActionID:  actionID,
		})
	}
}
func (this *OutcomeDetailParser) currentLineIsAction() bool {
	return this.isConcurrentTask() || this.isParallelTask()
}
func (this *OutcomeDetailParser) isConcurrentTask() bool {
	return this.words[0] == "-"
}
func (this *OutcomeDetailParser) isParallelTask() bool {
	first := this.words[0]
	first = strings.TrimRight(first, ".")
	number, _ := strconv.Atoi(first)
	return number > 0
}
func (this *OutcomeDetailParser) isCompletedTask() bool {
	return this.words[1] == "[X]" || this.words[1] == "[x]"
}
func (this *OutcomeDetailParser) isLatentTask() bool {
	return this.words[1] == "[?]"
}
func (this *OutcomeDetailParser) isIncompleteTask() bool {
	return this.words[1] == "[]" || (this.words[1] == "[" && this.words[2] == "]")
}
func (this *OutcomeDetailParser) parseActionDescription() string {
	start := strings.Index(this.line, "]") + 1
	return strings.TrimSpace(this.line[start:])
}
func (this *OutcomeDetailParser) parseOutcomeDescriptionLine() {
	_, _ = io.WriteString(this.description, this.line)
	_, _ = io.WriteString(this.description, "\n")
}

const trackOutcomeTemplate = `# {TITLE}

> {EXPLANATION}


## Actions:

-  [ ] concurrent @home
1. [ ] sequential @home


## Support Materials:`
