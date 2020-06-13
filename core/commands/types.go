package commands

type CreateResult struct {
	ID    string
	Error error
}

type Result struct {
	Error error
}

type Identifiable interface{ ID() string }

func (this TrackOutcome) ID() string                 { return this.Result.ID }
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
func (this ReorderActions) ID() string               { return this.OutcomeID }
func (this UpdateActionDescription) ID() string      { return this.OutcomeID }
func (this MarkActionStatusLatent) ID() string       { return this.OutcomeID }
func (this MarkActionStatusIncomplete) ID() string   { return this.OutcomeID }
func (this MarkActionStatusComplete) ID() string     { return this.OutcomeID }
func (this MarkActionStrategySequential) ID() string { return this.OutcomeID }
func (this MarkActionStrategyConcurrent) ID() string { return this.OutcomeID }
func (this DeleteAction) ID() string                 { return this.OutcomeID }

type Fallible interface{ Err() error }

func (this TrackOutcome) Err() error                 { return this.Result.Error }
func (this UpdateOutcomeTitle) Err() error           { return this.Result.Error }
func (this UpdateOutcomeExplanation) Err() error     { return this.Result.Error }
func (this UpdateOutcomeDescription) Err() error     { return this.Result.Error }
func (this DeleteOutcome) Err() error                { return this.Result.Error }
func (this DeclareOutcomeFixed) Err() error          { return this.Result.Error }
func (this DeclareOutcomeRealized) Err() error       { return this.Result.Error }
func (this DeclareOutcomeAbandoned) Err() error      { return this.Result.Error }
func (this DeclareOutcomeDeferred) Err() error       { return this.Result.Error }
func (this DeclareOutcomeUncertain) Err() error      { return this.Result.Error }
func (this TrackAction) Err() error                  { return this.Result.Error }
func (this ReorderActions) Err() error               { return this.Result.Error }
func (this UpdateActionDescription) Err() error      { return this.Result.Error }
func (this MarkActionStatusLatent) Err() error       { return this.Result.Error }
func (this MarkActionStatusIncomplete) Err() error   { return this.Result.Error }
func (this MarkActionStatusComplete) Err() error     { return this.Result.Error }
func (this MarkActionStrategySequential) Err() error { return this.Result.Error }
func (this MarkActionStrategyConcurrent) Err() error { return this.Result.Error }
func (this DeleteAction) Err() error                 { return this.Result.Error }
