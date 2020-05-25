package commands

type CreateResult struct {
	ID    string
	Error error
}

type Result struct {
	Error error
}

type Identifiable interface{ ID() string }

func (this UpdateOutcomeTitle) ID() string           { return this.OutcomeID }
func (this UpdateOutcomeExplanation) ID() string     { return this.OutcomeID }
func (this UpdateOutcomeDescription) ID() string     { return this.OutcomeID }
func (this DeleteOutcome) ID() string                { return this.OutcomeID }
func (this DeclareOutcomeFixed) ID() string          { return this.OutcomeID }
func (this DeclareOutcomeRealized) ID() string       { return this.OutcomeID }
func (this DeclareOutcomeAbandoned) ID() string      { return this.OutcomeID }
func (this DeclareOutcomeDeferred) ID() string       { return this.OutcomeID }
func (this DeclareOutcomeUncertain) ID() string      { return this.OutcomeID }
func (this TrackAction) ID() string                  { return this.OutcomeID }
func (this ChangeActionDisplayOrder) ID() string     { return this.OutcomeID }
func (this UpdateActionDescription) ID() string      { return this.OutcomeID }
func (this MarkActionStatusLatent) ID() string       { return this.OutcomeID }
func (this MarkActionStatusIncomplete) ID() string   { return this.OutcomeID }
func (this MarkActionStatusComplete) ID() string     { return this.OutcomeID }
func (this MarkActionStrategySequential) ID() string { return this.OutcomeID }
func (this MarkActionStrategyConcurrent) ID() string { return this.OutcomeID }
func (this DeleteAction) ID() string                 { return this.OutcomeID }
