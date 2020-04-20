package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		// Deprecated
		Stream chan interface{}
		Events []interface{}
	}
}
