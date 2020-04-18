package queries

type EventStream struct {
	Result struct {
		Stream chan interface{}
	}
}
