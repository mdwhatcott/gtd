package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Events []interface{}
	}
}
