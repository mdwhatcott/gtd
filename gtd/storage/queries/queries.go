package queries

type EventStream struct {
	UserID string
	Result struct {
		Stream chan interface{}
	}
}
