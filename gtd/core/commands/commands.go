package commands

type TrackOutcome struct {
	Title string

	Result struct {
		Error     error
		OutcomeID string
	}
}

type DefineOutcome struct {
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

	Result struct {
		Error error
	}
}

type UpdateActionDescription struct {
	OutcomeID      string
	ActionID       string
	NewDescription string

	Result struct {
		Error error
	}
}

type MarkActionStatusLatent struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionStatusIncomplete struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionStatusComplete struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionStrategySequential struct {
	OutcomeID string
	ActionID  string

	Result struct {
		Error error
	}
}

type MarkActionStrategyConcurrent struct {
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
