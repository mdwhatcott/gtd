package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Events []interface{}
	}
}

func (this OutcomeEventStream) ID() string {
	return this.OutcomeID
}

type EventStream struct {
	Result struct {
		Events []interface{}
	}
}
