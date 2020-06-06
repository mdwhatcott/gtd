package fake

import "fmt"

type Identifiable struct {
	Value interface{} `json:"value"`
}

func NewIdentifiable(_id interface{}) *Identifiable {
	return &Identifiable{Value: _id}
}

func (this *Identifiable) ID() string {
	return fmt.Sprint(this.Value)
}
