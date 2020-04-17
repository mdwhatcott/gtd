package commands

type TrackOutcome struct {
	UserID     string
	Definition string

	Result struct {
		Error     error
		OutcomeID string
	}
}

type RedefineOutcome struct {
	UserID        string
	OutcomeID     string
	NewDefinition string

	Result struct {
		Error error
	}
}

type DescribeOutcome struct {
	UserID      string
	OutcomeID   string
	Description string

	Result struct {
		Error error
	}
}

type DeleteOutcome struct {
	UserID    string
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeFixed struct {
	UserID    string
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeRealized struct {
	UserID    string
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeAbandoned struct {
	UserID    string
	OutcomeID string
	Reason    string

	Result struct {
		Error error
	}
}

type DeclareOutcomeDeferred struct {
	UserID    string
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeUncertain struct {
	UserID    string
	OutcomeID string

	Result struct {
		Error error
	}
}

type TrackAction struct {
	UserID     string
	OutcomeID  string
	Definition string
	Contexts   []string
	IsComplete bool
	Sequence   float64

	Result struct {
		Error    error
		ActionID string
	}
}

type ResequencedAction struct {
	UserID      string
	OutcomeID   string
	ActionID    string
	NewSequence float64

	Result struct {
		Error error
	}
}

type RedefineAction struct {
	UserID        string
	OutcomeID     string
	ActionID      string
	NewDefinition string

	Result struct {
		Error error
	}
}

type AddContextToAction struct {
	UserID     string
	OutcomeID  string
	ActionID   string
	NewContext string

	Result struct {
		Error error
	}
}

type RemoveContextFromAction struct {
	UserID         string
	OutcomeID      string
	ActionID       string
	RemovedContext string

	Result struct {
		Error error
	}
}

type MarkActionComplete struct {
	UserID    string
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionNotComplete struct {
	UserID    string
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type DeleteAction struct {
	UserID    string
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}
