package commands

type TrackOutcome struct {
	UserID     GUID
	Definition string

	Result struct {
		Error     error
		OutcomeID GUID
	}
}

type RedefineOutcome struct {
	UserID        GUID
	OutcomeID     GUID
	NewDefinition string

	Result struct {
		Error error
	}
}

type DescribeOutcome struct {
	UserID      GUID
	OutcomeID   GUID
	Description string

	Result struct {
		Error error
	}
}

type DeleteOutcome struct {
	UserID    GUID
	OutcomeID GUID

	Result struct {
		Error error
	}
}

type DeclareOutcomeFixed struct {
	UserID    GUID
	OutcomeID GUID

	Result struct {
		Error error
	}
}

type DeclareOutcomeRealized struct {
	UserID    GUID
	OutcomeID GUID

	Result struct {
		Error error
	}
}

type DeclareOutcomeAbandoned struct {
	UserID    GUID
	OutcomeID GUID
	Reason    string

	Result struct {
		Error error
	}
}

type DeclareOutcomeDeferred struct {
	UserID    GUID
	OutcomeID GUID

	Result struct {
		Error error
	}
}

type DeclareOutcomeUncertain struct {
	UserID    GUID
	OutcomeID GUID

	Result struct {
		Error error
	}
}

type TrackAction struct {
	UserID     GUID
	OutcomeID  GUID
	Definition string
	Contexts   []string
	IsComplete bool
	Sequence   float64

	Result struct {
		Error    error
		ActionID GUID
	}
}

type ResequencedAction struct {
	UserID      GUID
	OutcomeID   GUID
	ActionID    GUID
	NewSequence float64

	Result struct {
		Error error
	}
}

type RedefinedAction struct {
	UserID        GUID
	OutcomeID     GUID
	ActionID      GUID
	NewDefinition string

	Result struct {
		Error error
	}
}

type AddContextToAction struct {
	UserID     GUID
	OutcomeID  GUID
	ActionID   GUID
	NewContext string

	Result struct {
		Error error
	}
}

type RemoveContextFromAction struct {
	UserID         GUID
	OutcomeID      GUID
	ActionID       GUID
	RemovedContext string

	Result struct {
		Error error
	}
}

type MarkActionComplete struct {
	UserID    GUID
	OutcomeID GUID
	ActionID  GUID

	Result struct {
		Error error
	}
}

type MarkActionNotComplete struct {
	UserID    GUID
	OutcomeID GUID
	ActionID  GUID

	Result struct {
		Error error
	}
}

type DeleteAction struct {
	UserID    GUID
	OutcomeID GUID
	ActionID  GUID

	Result struct {
		Error error
	}
}
