package events

type OutcomeTrackedV1 struct {
	Timestamp Time
	OutcomeID string
	Title     string
}
type OutcomeTitleUpdatedV1 struct {
	Timestamp    Time
	OutcomeID    string
	UpdatedTitle string
}
type OutcomeExplanationUpdatedV1 struct {
	Timestamp          Time
	OutcomeID          string
	UpdatedExplanation string
}
type OutcomeDescriptionUpdatedV1 struct {
	Timestamp          Time
	OutcomeID          string
	UpdatedDescription string
}
type OutcomeDeletedV1 struct {
	Timestamp Time
	OutcomeID string
}
type OutcomeFixedV1 struct {
	Timestamp Time
	OutcomeID string
}
type OutcomeRealizedV1 struct {
	Timestamp Time
	OutcomeID string
}
type OutcomeAbandonedV1 struct {
	Timestamp Time
	OutcomeID string
}
type OutcomeDeferredV1 struct {
	Timestamp Time
	OutcomeID string
}
type OutcomeUncertainV1 struct {
	Timestamp Time
	OutcomeID string
}
type ActionTrackedV1 struct {
	Timestamp   Time
	OutcomeID   string
	ActionID    string
	Description string
	Contexts    []string
}
type ActionsReorderedV1 struct {
	Timestamp    Time
	OutcomeID    string
	ReorderedIDs []string
}
type ActionDescriptionUpdatedV1 struct {
	Timestamp          Time
	OutcomeID          string
	ActionID           string
	UpdatedDescription string
	UpdatedContexts    []string
}
type ActionStatusMarkedLatentV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
type ActionStatusMarkedIncompleteV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
type ActionStatusMarkedCompleteV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
type ActionStrategyMarkedSequentialV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
type ActionStrategyMarkedConcurrentV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
type ActionDeletedV1 struct {
	Timestamp Time
	OutcomeID string
	ActionID  string
}
