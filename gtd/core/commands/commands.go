package commands

type TrackOutcome struct {
	Definition string

	Result struct {
		Error     error
		OutcomeID string
	}
}

type RedefineOutcome struct {
	OutcomeID     string
	NewDefinition string

	Result struct {
		Error error
	}
}

type DescribeOutcome struct {
	OutcomeID   string
	Description string

	Result struct {
		Error error
	}
}

type DeleteOutcome struct {
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeFixed struct {
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeRealized struct {
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeAbandoned struct {
	OutcomeID string
	Reason    string

	Result struct {
		Error error
	}
}

type DeclareOutcomeDeferred struct {
	OutcomeID string

	Result struct {
		Error error
	}
}

type DeclareOutcomeUncertain struct {
	OutcomeID string

	Result struct {
		Error error
	}
}

type TrackAction struct {
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
	OutcomeID   string
	ActionID    string
	NewSequence float64

	Result struct {
		Error error
	}
}

type RedefineAction struct {
	OutcomeID     string
	ActionID      string
	NewDefinition string

	Result struct {
		Error error
	}
}

type AddContextToAction struct {
	OutcomeID  string
	ActionID   string
	NewContext string

	Result struct {
		Error error
	}
}

type RemoveContextFromAction struct {
	OutcomeID      string
	ActionID       string
	RemovedContext string

	Result struct {
		Error error
	}
}

type MarkActionComplete struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionNotComplete struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type DeleteAction struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}
