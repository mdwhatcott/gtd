package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Events []interface{}
	}
}

type EventStream struct {
	Result struct {
		Events []interface{}
	}
}
