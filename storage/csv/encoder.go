package csv

import (
	"encoding/csv"
	"io"
	"reflect"
)

type Encoder struct {
	writer   *csv.Writer
	registry map[string]func(interface{}) []string
}

func NewEncoder(writer_ io.Writer, registry_ map[string]func(interface{}) []string) *Encoder {
	WRITER := csv.NewWriter(writer_)
	WRITER.Comma = '\t'
	return &Encoder{
		writer:   WRITER,
		registry: registry_,
	}
}

func (this *Encoder) Encode(v interface{}) (err_ error) {
	defer this.writer.Flush()
	TYPE := reflect.TypeOf(v).String()
	RECORD := this.registry[TYPE](v)
	return this.writer.Write(RECORD)
}
