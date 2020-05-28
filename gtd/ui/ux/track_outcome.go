package ux

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/ui"
)

type TrackOutcomeExperience struct {
	handler core.Handler
	scanner *bufio.Scanner
	line    string
	words   []string
	parse   func()

	outcome     *commands.TrackOutcome
	description *strings.Builder
}

func NewTrackOutcomeExperience(handler core.Handler, editor ui.Editor) *TrackOutcomeExperience {
	content := editor.EditTempFile(trackOutcomeTemplate)
	scanner := bufio.NewScanner(strings.NewReader(content))
	ux := &TrackOutcomeExperience{
		handler:     handler,
		scanner:     scanner,
		description: new(strings.Builder),
	}
	ux.parse = ux.parseTitleLine
	return ux
}

func (this *TrackOutcomeExperience) Engage() error {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		this.words = strings.Fields(this.line)
		this.parse()
		if this.outcome == nil {
			break
		}
	}
	if this.description.Len() > 0 {
		this.handler.Handle(&commands.UpdateOutcomeDescription{
			OutcomeID:          this.outcome.Result.ID,
			UpdatedDescription: this.description.String(),
		})
	}
	return nil
}

func (this *TrackOutcomeExperience) parseTitleLine() {
	if !strings.HasPrefix(this.line, "# ") {
		return
	}
	if this.line == "# {TITLE}" {
		return
	}
	this.outcome = &commands.TrackOutcome{Title: this.line[2:]}
	this.handler.Handle(this.outcome)
	this.parse = this.parseExplanationLine
}
func (this *TrackOutcomeExperience) parseExplanationLine() {
	if this.line == "## Actions:" {
		this.parse = this.parseActionLine
		return
	}
	if this.line == "" {
		return
	}
	if strings.HasPrefix(this.line, "> ") {
		this.handler.Handle(&commands.UpdateOutcomeExplanation{
			OutcomeID:          this.outcome.Result.ID,
			UpdatedExplanation: this.line[2:],
		})
	}
}
func (this *TrackOutcomeExperience) parseActionLine() {
	if this.line == "## Support Materials:" {
		this.parse = this.parseOutcomeDescriptionLine
	}
	if this.line == "" {
		return
	}
	if !this.currentLineIsAction() {
		return
	}

	this.parseAction()
}
func (this *TrackOutcomeExperience) parseAction() {
	action := &commands.TrackAction{
		OutcomeID:   this.outcome.Result.ID,
		Description: this.parseActionDescription(),
	}
	this.handler.Handle(action)

	this.parseActionStrategy(action.Result.ID)
	this.parseActionStatus(action.Result.ID)
}
func (this *TrackOutcomeExperience) parseActionStatus(actionID string) {
	if this.isCompletedTask() {
		this.handler.Handle(&commands.MarkActionStatusComplete{
			OutcomeID: this.outcome.Result.ID,
			ActionID:  actionID,
		})
	} else if this.isIncompleteTask() {
		this.handler.Handle(&commands.MarkActionStatusIncomplete{
			OutcomeID: this.outcome.Result.ID,
			ActionID:  actionID,
		})
	} else if this.isLatentTask() {
		this.handler.Handle(&commands.MarkActionStatusLatent{
			OutcomeID: this.outcome.Result.ID,
			ActionID:  actionID,
		})
	}
}
func (this *TrackOutcomeExperience) parseActionStrategy(actionID string) {
	if this.isConcurrentTask() {
		this.handler.Handle(&commands.MarkActionStrategyConcurrent{
			OutcomeID: this.outcome.Result.ID,
			ActionID:  actionID,
		})
	} else if this.isParallelTask() {
		this.handler.Handle(&commands.MarkActionStrategySequential{
			OutcomeID: this.outcome.Result.ID,
			ActionID:  actionID,
		})
	}
}
func (this *TrackOutcomeExperience) currentLineIsAction() bool {
	return this.isConcurrentTask() || this.isParallelTask()
}
func (this *TrackOutcomeExperience) isConcurrentTask() bool {
	return this.words[0] == "-"
}
func (this *TrackOutcomeExperience) isParallelTask() bool {
	first := this.words[0]
	first = strings.TrimRight(first, ".")
	number, _ := strconv.Atoi(first)
	return number > 0
}
func (this *TrackOutcomeExperience) isCompletedTask() bool {
	return this.words[1] == "[X]" || this.words[1] == "[x]"
}
func (this *TrackOutcomeExperience) isLatentTask() bool {
	return this.words[1] == "[?]"
}
func (this *TrackOutcomeExperience) isIncompleteTask() bool {
	return this.words[1] == "[]" || (this.words[1] == "[" && this.words[2] == "]")
}
func (this *TrackOutcomeExperience) parseActionDescription() string {
	start := strings.Index(this.line, "]") + 1
	return strings.TrimSpace(this.line[start:])
}
func (this *TrackOutcomeExperience) parseOutcomeDescriptionLine() {
	_, _ = io.WriteString(this.description, this.line)
	_, _ = io.WriteString(this.description, "\n")
}

const trackOutcomeTemplate = `# {TITLE}

> {EXPLANATION}


## Actions:

-  [ ] concurrent @home
1. [ ] sequential @home


## Support Materials:`
