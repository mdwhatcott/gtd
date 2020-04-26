package json

import (
	"bytes"
	"io"
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

	source   *bytes.Buffer
	decoder  *Decoder
	registry map[string]reflect.Type
}

func (this *DecoderFixture) Setup() {
	this.source = new(bytes.Buffer)
	this.decoder = NewDecoder(this.source)
	this.registry = map[string]reflect.Type{
		"string": reflect.TypeOf(""),
	}
}

func (this *DecoderFixture) TestDecodeAll() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"b"` + "\n")

	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"c"` + "\n")

	all := this.decoder.DecodeAll(this.registry)

	this.So(gather(all), should.Resemble, []interface{}{"a", "b", "c"})
}

func (this *DecoderFixture) TestDecodeAll_MissingValue() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"b"` + "\n")

	this.source.WriteString(`"string"` + "\n")
	//this.source.WriteString(`"c"` + "\n") // missing last line!

	all := this.decoder.DecodeAll(this.registry)

	this.So(gather(all), should.Resemble, []interface{}{"a", "b"})
	this.So(this.decoder.Error(), should.Wrap, io.EOF)
}
func (this *DecoderFixture) TestDecodeAll_UnregisteredType() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"b"` + "\n")

	this.source.WriteString(`"UNREGISTERED-TYPE"` + "\n")
	this.source.WriteString(`"c"` + "\n")

	all := this.decoder.DecodeAll(this.registry)

	this.So(gather(all), should.Resemble, []interface{}{"a", "b"})
	this.So(this.decoder.Error(), should.Wrap, errUnregisteredType)
	this.So(this.decoder.Error().Error(), should.ContainSubstring, "UNREGISTERED-TYPE")
}
func gather(queue chan interface{}) (all []interface{}) {
	for item := range queue {
		all = append(all, item)
	}
	return all
}
