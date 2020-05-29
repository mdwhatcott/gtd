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
type DeclareOutcomeUncertain struct {
	OutcomeID string
	Result    Result
}
type TrackAction struct {
	OutcomeID   string
	Description string
	Result      CreateResult
}
type UpdateActionDescription struct {
	OutcomeID      string
	ActionID       string
	NewDescription string // TODO: rename: UpdatedDescription
	Result         Result
}
type ReorderActions struct {
	OutcomeID  string
	NewIDOrder []string
	Result     Result
}
type MarkActionStatusLatent struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusIncomplete struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStatusComplete struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategySequential struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
type MarkActionStrategyConcurrent struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
type DeleteAction struct {
	OutcomeID string
	ActionID  string
	Result    Result
}
