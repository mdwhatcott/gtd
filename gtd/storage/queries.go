package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Stream chan interface{}
	}
}
