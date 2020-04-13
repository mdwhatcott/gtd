package events

import "time"

type OutcomeIdentifiedV1 struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	OutcomeID  string    `json:"outcome_id"`
	Definition string    `json:"definition"`
}

type OutcomeRedefinedV1 struct {
	Timestamp     time.Time `json:"timestamp"`
	UserID        string    `json:"user_id"`
	OutcomeID     string    `json:"outcome_id"`
	NewDefinition string    `json:"new_definition"`
}

type OutcomeDescribedV1 struct {
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id"`
	OutcomeID   string    `json:"outcome_id"`
	Description string    `json:"description"`
}

type OutcomeDeletedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
}

type OutcomeIsCurrentV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
}

type OutcomeIsCompletedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
}

type OutcomeIsRejectedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
	Reason    string    `json:"reason"`
}

type OutcomeIsASomedayV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
}

type OutcomeIsAMaybeV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
}

type ActionIdentifiedV1 struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	OutcomeID  string    `json:"outcome_id"`
	ActionID   string    `json:"action_id"`
	Definition string    `json:"definition"`
	Contexts   []string  `json:"contexts"`
	IsComplete bool      `json:"is_complete"`
	Sequence   float64   `json:"sequence"`
}

type ActionResequencedV1 struct {
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id"`
	OutcomeID   string    `json:"outcome_id"`
	ActionID    string    `json:"action_id"`
	NewSequence float64   `json:"new_definition"`
}

type ActionRedefinedV1 struct {
	Timestamp     time.Time `json:"timestamp"`
	UserID        string    `json:"user_id"`
	OutcomeID     string    `json:"outcome_id"`
	ActionID      string    `json:"action_id"`
	NewDefinition string    `json:"new_definition"`
}

type ActionContextAddedV1 struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	OutcomeID  string    `json:"outcome_id"`
	ActionID   string    `json:"action_id"`
	NewContext string    `json:"new_context"`
}

type ActionContextRemovedV1 struct {
	Timestamp      time.Time `json:"timestamp"`
	UserID         string    `json:"user_id"`
	OutcomeID      string    `json:"outcome_id"`
	ActionID       string    `json:"action_id"`
	RemovedContext string    `json:"removed_context"`
}

type ActionCompletedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
	ActionID  string    `json:"action_id"`
}

type ActionUncompletedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
	ActionID  string    `json:"action_id"`
}

type ActionDeletedV1 struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	OutcomeID string    `json:"outcome_id"`
	ActionID  string    `json:"action_id"`
}
