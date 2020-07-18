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

func NewEncoder(writer io.Writer, registry map[string]func(interface{}) []string) *Encoder {
	WRITER := csv.NewWriter(writer)
	WRITER.Comma = '\t'
	return &Encoder{
		writer:   WRITER,
		registry: registry,
	}
}

func (this *Encoder) Encode(v interface{}) (err_ error) {
	defer this.writer.Flush()
	TYPE := reflect.TypeOf(v).String()
	RECORD := this.registry[TYPE](v)
	return this.writer.Write(RECORD)
}
