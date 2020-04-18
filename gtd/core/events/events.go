package events

type OutcomeTrackedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	Title     string `json:"definition"`
}

type OutcomeExplanationProvidedV1 struct {
	Timestamp   Time   `json:"timestamp"`
	OutcomeID   string `json:"outcome_id"`
	Explanation string `json:"explanation"`
}

type OutcomeDescriptionUpdatedV1 struct {
	Timestamp      Time   `json:"timestamp"`
	OutcomeID      string `json:"outcome_id"`
	NewDescription string `json:"description"`
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
	Reason    string `json:"reason"`
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
	IsComplete  bool     `json:"is_complete"`
	Sequence    float64  `json:"sequence"`
}

type ActionReorderedV1 struct {
	Timestamp   Time    `json:"timestamp"`
	OutcomeID   string  `json:"outcome_id"`
	ActionID    string  `json:"action_id"`
	NewSequence float64 `json:"new_definition"`
}

type ActionDescriptionUpdatedV1 struct {
	Timestamp     Time   `json:"timestamp"`
	OutcomeID     string `json:"outcome_id"`
	ActionID      string `json:"action_id"`
	NewDefinition string `json:"new_definition"`
}

type ActionStatusMarkedLatentV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionStatusMarkedIncompleteV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionStatusMarkedCompleteV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionStrategyMarkedSequentialV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionStrategyMarkedConcurrentV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionDeletedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
