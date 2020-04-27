package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Events []interface{} // TODO: go back to a chan interface{}
	}
}
