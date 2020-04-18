package events

type OutcomeDefinedV1 struct {
	Timestamp  Time   `json:"timestamp"`
	OutcomeID  string `json:"outcome_id"`
	Definition string `json:"definition"`
}

type OutcomeRedefinedV1 struct {
	Timestamp     Time   `json:"timestamp"`
	OutcomeID     string `json:"outcome_id"`
	NewDefinition string `json:"new_definition"`
}

type OutcomeDescribedV1 struct {
	Timestamp   Time   `json:"timestamp"`
	OutcomeID   string `json:"outcome_id"`
	Description string `json:"description"`
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

type ActionIdentifiedV1 struct {
	Timestamp  Time     `json:"timestamp"`
	OutcomeID  string   `json:"outcome_id"`
	ActionID   string   `json:"action_id"`
	Definition string   `json:"definition"`
	Contexts   []string `json:"contexts"`
	IsComplete bool     `json:"is_complete"`
	Sequence   float64  `json:"sequence"`
}

type ActionResequencedV1 struct {
	Timestamp   Time    `json:"timestamp"`
	OutcomeID   string  `json:"outcome_id"`
	ActionID    string  `json:"action_id"`
	NewSequence float64 `json:"new_definition"`
}

type ActionRedefinedV1 struct {
	Timestamp     Time   `json:"timestamp"`
	OutcomeID     string `json:"outcome_id"`
	ActionID      string `json:"action_id"`
	NewDefinition string `json:"new_definition"`
}

type ActionContextAddedV1 struct {
	Timestamp  Time   `json:"timestamp"`
	OutcomeID  string `json:"outcome_id"`
	ActionID   string `json:"action_id"`
	NewContext string `json:"new_context"`
}

type ActionContextRemovedV1 struct {
	Timestamp      Time   `json:"timestamp"`
	OutcomeID      string `json:"outcome_id"`
	ActionID       string `json:"action_id"`
	RemovedContext string `json:"removed_context"`
}

type ActionMarkedCompleteV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionMarkedNotCompleteV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}

type ActionDeletedV1 struct {
	Timestamp Time   `json:"timestamp"`
	OutcomeID string `json:"outcome_id"`
	ActionID  string `json:"action_id"`
}
