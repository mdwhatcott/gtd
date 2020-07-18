package fake

import "fmt"

type Identifiable struct {
	Value interface{} `json:"value"`
}

func NewIdentifiable(id interface{}) *Identifiable {
	return &Identifiable{Value: id}
}

func (this *Identifiable) ID() string {
	return fmt.Sprint(this.Value)
}
