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
type DeleteOutcome struct {
	OutcomeID string
	Result    Result
}
type DeclareOutcomeFixed struct {
	OutcomeID string
	Result    Result
}
type DeclareOutcomeRealized struct {
	OutcomeID string
	Result    Result
}
type DeclareOutcomeAbandoned struct {
	OutcomeID string
	Result    Result
}
type DeclareOutcomeDeferred struct {
	OutcomeID string
	Result    Result
}
type DeclareOutcomeUncertain struct { // joyride
	OutcomeID string
	Result    Result
}
type TrackAction struct { // joyride
	OutcomeID   string
	Description string
	Result      CreateResult
}
type ReorderAction struct { // joyride
	OutcomeID string
	ActionID  string
	NewOrder  float64
	Result    Result
}
type UpdateActionDescription struct { // joyride
	OutcomeID      string
	ActionID       string
	NewDescription string
	Result         Result
}
type MarkActionStatusLatent struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusIncomplete struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusComplete struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategySequential struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategyConcurrent struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
type DeleteAction struct { // joyride
	OutcomeID string
	ActionID  string
	Result    Result
}
