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

	source   *bytes.Buffer
	decoder  *Decoder
	registry map[string]reflect.Type
}

func (this *DecoderFixture) Setup() {
	this.source = new(bytes.Buffer)
	this.registry = map[string]reflect.Type{"string": reflect.TypeOf("s")}
	this.decoder = NewDecoder(this.source, this.registry)
}

func (this *DecoderFixture) TestDecodeType() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.BeNil)
	this.So(V, should.Equal, "a")
}

func (this *DecoderFixture) TestDecode_FailureToDecodeTypeName_Err() {
	this.source.WriteString(`invalid json`)

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingInvalidJSONTypeName)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_FailureToDecodeValue_Err() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`invalid json`)

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingInvalidJSONValue)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_UnrecognizedType_Err() {
	this.source.WriteString(`"UNRECOGNIZED_TYPE"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingUnregisteredType)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_TypeValueMismatch_Err() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`1` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingTypeMismatch)
	this.So(V, should.BeNil)
}
