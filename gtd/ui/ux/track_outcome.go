package ux

import (
	"bufio"
	"strings"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/commands"
	"github.com/mdwhatcott/gtd/gtd/ui"
)

type TrackOutcomeExperience struct {
	handler core.Handler
	scanner *bufio.Scanner
	line    string
	parse   func()

	outcome *commands.TrackOutcome
}

func NewTrackOutcomeExperience(handler core.Handler, editor ui.Editor) *TrackOutcomeExperience {
	content := editor.EditTempFile(trackOutcomeTemplate)
	scanner := bufio.NewScanner(strings.NewReader(content))
	ux := &TrackOutcomeExperience{
		handler: handler,
		scanner: scanner,
	}
	ux.parse = ux.parseTitle
	return ux
}

func (this *TrackOutcomeExperience) Engage() error {
	for this.scanner.Scan() {
		this.line = this.scanner.Text()
		this.parse()
		if this.outcome == nil {
			break
		}
	}
	return nil
}

func (this *TrackOutcomeExperience) parseTitle() {
	if !strings.HasPrefix(this.line, "# ") {
		return
	}
	if this.line == "# {TITLE}" {
		return
	}
	this.outcome = &commands.TrackOutcome{Title: this.line[2:]}
	this.handler.Handle(this.outcome)
	this.parse = this.parseExplanation
}
func (this *TrackOutcomeExperience) parseExplanation() {
	if this.line == "## Actions:" {
		this.parse = this.parseAction
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
func (this *TrackOutcomeExperience) parseAction() {
	if this.line == "## Support Materials:" {
		this.parse = this.parseDescription
	}
	if this.line == "" {
		return
	}
	if strings.HasPrefix(this.line, "-  [") || strings.HasPrefix(this.line, "1. [") {
		this.handler.Handle(&commands.TrackAction{
			OutcomeID:   this.outcome.Result.ID,
			Description: this.line[7:],
		})
	}
}
func (this *TrackOutcomeExperience) parseDescription() {}

const trackOutcomeTemplate = `# {TITLE}

> {EXPLANATION}


## Actions:

-  [ ] concurrent @home
1. [ ] sequential @home


## Support Materials:`
