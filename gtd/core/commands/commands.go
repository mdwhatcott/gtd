package commands

type TrackOutcome struct {
	Title string

	Result CreateResult
}

type UpdateOutcomeTitle struct {
	OutcomeID string
	NewTitle  string

	Result Result
}

type UpdateOutcomeExplanation struct {
	OutcomeID   string
	Explanation string

	Result Result
}

type UpdateOutcomeDescription struct {
	OutcomeID          string
	UpdatedDescription string

	Result Result
}

type DeleteOutcome struct {
	OutcomeID string

	Result Result
}

type DeclareOutcomeFixed struct {
	OutcomeID string

	Result Result
}

type DeclareOutcomeRealized struct {
	OutcomeID string

	Result Result
}

type DeclareOutcomeAbandoned struct {
	OutcomeID string
	Reason    string

	Result Result
}

type DeclareOutcomeDeferred struct {
	OutcomeID string

	Result Result
}

type DeclareOutcomeUncertain struct {
	OutcomeID string

	Result Result
}

type TrackAction struct {
	OutcomeID   string
	Description string

	Result struct {
		Error    error
		ActionID string
	}
}

type ReorderAction struct {
	OutcomeID string
	ActionID  string
	NewOrder  float64

	Result Result
}

type UpdateActionDescription struct {
	OutcomeID      string
	ActionID       string
	NewDescription string

	Result Result
}

type MarkActionStatusLatent struct {
	OutcomeID string
	ActionID  string

	Result Result
}

type MarkActionStatusIncomplete struct {
	OutcomeID string
	ActionID  string

	Result Result
}

type MarkActionStatusComplete struct {
	OutcomeID string
	ActionID  string

	Result Result
}

type MarkActionStrategySequential struct {
	OutcomeID string
	ActionID  string

	Result Result
}

type MarkActionStrategyConcurrent struct {
	OutcomeID string
	ActionID  string

	Result Result
}

type DeleteAction struct {
	OutcomeID string
	ActionID  string

	Result Result
}
