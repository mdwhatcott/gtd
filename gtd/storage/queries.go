package storage

type OutcomeEventStream struct {
	OutcomeID string

	Result struct {
		Stream chan interface{}
	}
}

func (this *OutcomeEventStream) ID() string { return this.OutcomeID }
