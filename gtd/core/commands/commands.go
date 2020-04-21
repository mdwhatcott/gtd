package commands

type TrackOutcome struct {
	Title  string
	Result CreateResult
}
type UpdateOutcomeTitle struct {
	OutcomeID    string
	UpdatedTitle string
	Result       Result
}
type UpdateOutcomeExplanation struct {
	OutcomeID          string
	UpdatedExplanation string
	Result             Result
}
type UpdateOutcomeDescription struct {
	OutcomeID          string
	UpdatedDescription string
	Result             Result
}
type DeleteOutcome struct { // todo: joyride
	OutcomeID string
	Result    Result
}
type DeclareOutcomeFixed struct { // todo: joyride
	OutcomeID string
	Result    Result
}
type DeclareOutcomeRealized struct { // todo: joyride
	OutcomeID string
	Result    Result
}
type DeclareOutcomeAbandoned struct { // todo: joyride
	OutcomeID string
	Reason    string
	Result    Result
}
type DeclareOutcomeDeferred struct { // todo: joyride
	OutcomeID string
	Result    Result
}
type DeclareOutcomeUncertain struct { // todo: joyride
	OutcomeID string
	Result    Result
}
type TrackAction struct { // todo: joyride
	OutcomeID   string
	Description string
	Result      struct {
		Error    error
		ActionID string
	}
}
type ReorderAction struct { // todo: joyride
	OutcomeID string
	ActionID  string
	NewOrder  float64
	Result    Result
}
type UpdateActionDescription struct { // todo: joyride
	OutcomeID      string
	ActionID       string
	NewDescription string
	Result         Result
}
type MarkActionStatusLatent struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusIncomplete struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusComplete struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategySequential struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategyConcurrent struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type DeleteAction struct { // todo: joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
