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

	v, err := this.decoder.Decode()

	this.So(err, should.BeNil)
	this.So(v, should.Equal, "a")
}

func (this *DecoderFixture) TestDecode_FailureToDecodeTypeName_Err() {
	this.source.WriteString(`invalid json`)

	v, err := this.decoder.Decode()

	this.So(err, should.Wrap, ErrDecodingInvalidJSONTypeName)
	this.So(v, should.BeNil)
}

func (this *DecoderFixture) TestDecode_FailureToDecodeValue_Err() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`invalid json`)

	v, err := this.decoder.Decode()

	this.So(err, should.Wrap, ErrDecodingInvalidJSONValue)
	this.So(v, should.BeNil)
}

func (this *DecoderFixture) TestDecode_UnrecognizedType_Err() {
	this.source.WriteString(`"UNRECOGNIZED_TYPE"` + "\n")
	this.source.WriteString(`"a"` + "\n")

	v, err := this.decoder.Decode()

	this.So(err, should.Wrap, ErrDecodingUnregisteredType)
	this.So(v, should.BeNil)
}

func (this *DecoderFixture) TestDecode_TypeValueMismatch_Err() {
	this.source.WriteString(`"string"` + "\n")
	this.source.WriteString(`1` + "\n")

	v, err := this.decoder.Decode()

	this.So(err, should.Wrap, ErrDecodingTypeMismatch)
	this.So(v, should.BeNil)
}
