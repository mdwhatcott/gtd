package events

type OutcomeTrackedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	Title     string `json:"title,omitempty"`
}
type OutcomeTitleUpdatedV1 struct {
	Timestamp    Time   `json:"timestamp,omitempty"`
	OutcomeID    string `json:"outcome_id,omitempty"`
	UpdatedTitle string `json:"updated_title,omitempty"`
}
type OutcomeExplanationUpdatedV1 struct {
	Timestamp          Time   `json:"timestamp,omitempty"`
	OutcomeID          string `json:"outcome_id,omitempty"`
	UpdatedExplanation string `json:"explanation,omitempty"`
}
type OutcomeDescriptionUpdatedV1 struct {
	Timestamp          Time   `json:"timestamp,omitempty"`
	OutcomeID          string `json:"outcome_id,omitempty"`
	UpdatedDescription string `json:"description,omitempty"`
}
type OutcomeDeletedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type OutcomeFixedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type OutcomeRealizedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type OutcomeAbandonedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type OutcomeDeferredV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type OutcomeUncertainV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
}
type ActionTrackedV1 struct {
	Timestamp   Time     `json:"timestamp,omitempty"`
	OutcomeID   string   `json:"outcome_id,omitempty"`
	ActionID    string   `json:"action_id,omitempty"`
	Description string   `json:"definition,omitempty"`
	Contexts    []string `json:"contexts,omitempty"`
}
type ActionsReorderedV1 struct {
	Timestamp    Time     `json:"timestamp,omitempty"`
	OutcomeID    string   `json:"outcome_id,omitempty"`
	ReorderedIDs []string `json:"reordered_ids,omitempty"`
}
type ActionDescriptionUpdatedV1 struct {
	Timestamp          Time     `json:"timestamp,omitempty"`
	OutcomeID          string   `json:"outcome_id,omitempty"`
	ActionID           string   `json:"action_id,omitempty"`
	UpdatedDescription string   `json:"updated_definition,omitempty"`
	UpdatedContexts    []string `json:"updated_contexts,omitempty"`
}
type ActionStatusMarkedLatentV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
type ActionStatusMarkedIncompleteV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
type ActionStatusMarkedCompleteV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
type ActionStrategyMarkedSequentialV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
type ActionStrategyMarkedConcurrentV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
type ActionDeletedV1 struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	OutcomeID string `json:"outcome_id,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
}
