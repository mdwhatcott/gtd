package events

type OutcomeTrackedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	Title     string `json:"title"`
}
type OutcomeTitleUpdatedV1 struct {
	Timestamp    Time   `json:"timestamp"`
	OutcomeID    string `json:"outcome_id"`
	UpdatedTitle string `json:"new_title"`
}
type OutcomeExplanationUpdatedV1 struct {
	Timestamp          Time   `json:"timestamp"`
	OutcomeID          string `json:"outcome_id"`
	UpdatedExplanation string `json:"explanation"`
}
type OutcomeDescriptionUpdatedV1 struct {
	Timestamp          Time   `json:"timestamp"`
	OutcomeID          string `json:"outcome_id"`
	UpdatedDescription string `json:"description"`
}
type OutcomeDeletedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type OutcomeFixedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type OutcomeRealizedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type OutcomeAbandonedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type OutcomeDeferredV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type OutcomeUncertainV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
}
type ActionTrackedV1 struct {
	Timestamp   Time     `json:"timestamp"`
	OutcomeID   string   `json:"outcome_id"`
	ActionID    string   `json:"action_id"`
	Description string   `json:"definition"`
	Contexts    []string `json:"contexts"`
	Sequence    int      `json:"sequence"`
}
type ActionsReorderedV1 struct { // TODO: outcome details
	Timestamp  Time     `json:"timestamp"`
	OutcomeID  string   `json:"outcome_id"`
	NewIDOrder []string `json:"new_id_order"`
}
type ActionDescriptionUpdatedV1 struct { // TODO: outcome details
	Timestamp          Time     `json:"timestamp"`
	OutcomeID          string   `json:"outcome_id"`
	ActionID           string   `json:"action_id"`
	UpdatedDescription string   `json:"updated_definition"`
	UpdatedContexts    []string `json:"updated_contexts"`
}
type ActionStatusMarkedLatentV1 struct { // TODO: outcome details
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
type ActionStatusMarkedIncompleteV1 struct { // TODO: outcome details
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
type ActionStatusMarkedCompleteV1 struct { // TODO: outcome details
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
type ActionStrategyMarkedSequentialV1 struct { // TODO: outcome details
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
type ActionStrategyMarkedConcurrentV1 struct { // TODO: outcome details
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
type ActionDeletedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
