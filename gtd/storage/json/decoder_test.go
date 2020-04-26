package json

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDecoderFixture(t *testing.T) {
	gunit.Run(new(DecoderFixture), t)
}

type DecoderFixture struct {
	*gunit.Fixture

	source  *bytes.Buffer
	decoder *Decoder
}

func (this *DecoderFixture) Setup() {
	this.source = new(bytes.Buffer)
	this.decoder = NewDecoder(this.source)
}

func (this *DecoderFixture) Test() {
	this.source.WriteString(`"int"` + "\n")
	this.source.WriteString("42" + "\n")
	registry := map[string]reflect.Type{
		"int": reflect.TypeOf(42),
	}

	value := this.decoder.DecodeType(registry)
	err := this.decoder.Decode(&value)

	this.So(err, should.BeNil)
	this.So(value, should.Equal, 42)
}
