package csv

import (
	"encoding/csv"
	"io"
)

type Decoder struct {
	reader   *csv.Reader
	registry map[string]func([]string) interface{}
}

func NewDecoder(reader_ io.Reader, registry_ map[string]func([]string) interface{}) *Decoder {
	READER := csv.NewReader(reader_)
	READER.Comma = '\t'
	READER.FieldsPerRecord = -1
	return &Decoder{
		reader:   READER,
		registry: registry_,
	}
}

func (this *Decoder) Decode() (interface{}, error) {
	RECORD, ERR := this.reader.Read()
	if ERR != nil {
		return nil, ERR
	}
	return this.registry[RECORD[2]](RECORD), ERR
}
