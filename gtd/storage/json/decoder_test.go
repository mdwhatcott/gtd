package json

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/gtd/core"
)

func TestDecoderFixture(t *testing.T) {
	gunit.Run(new(DecoderFixture), t)
}

type DecoderFixture struct {
	*gunit.Fixture

	source   *bytes.Buffer
	decoder  *Decoder
	registry map[string]core.Transformer
}

func (this *DecoderFixture) Setup() {
	this.source = new(bytes.Buffer)
	this.registry = map[string]core.Transformer{
		reflect.TypeOf(Struct{}).String(): TransformStruct,
	}
	this.decoder = NewDecoder(this.source, this.registry)
}

func (this *DecoderFixture) TestDecodeStructType() {
	this.Println(this.registry)
	this.source.WriteString(`"json.Struct"` + "\n")
	this.source.WriteString(`{"message":"Hello, world!"}` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.BeNil)
	this.So(V, should.Resemble, Struct{Message: "Hello, world!"})
}

func (this *DecoderFixture) TestDecode_EOF() {
	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Equal, io.EOF)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_FailureToDecodeTypeName_Err() {
	this.source.WriteString(`invalid json`)

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingInvalidJSONTypeName)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_FailureToDecodeValue_Err() {
	this.source.WriteString(`"json.Struct"` + "\n")
	this.source.WriteString(`invalid json`)

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingInvalidJSONValue)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_UnrecognizedType_Err() {
	this.source.WriteString(`"UNRECOGNIZED_TYPE"` + "\n")
	this.source.WriteString(`{"message":"Hello, world!"}` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingUnregisteredType)
	this.So(V, should.BeNil)
}

func (this *DecoderFixture) TestDecode_TypeValueMismatch_Err() {
	this.source.WriteString(`"json.Struct"` + "\n")
	this.source.WriteString(`1` + "\n")

	V, ERR := this.decoder.Decode()

	this.So(ERR, should.Wrap, ErrDecodingInvalidJSONValue)
	this.So(V, should.BeNil)
}

////////////////////////////////////////////////////////////////

func TransformStruct(raw map[string]interface{}) interface{} {
	return Struct{
		Message: raw["message"].(string),
	}
}

type Struct struct {
	Message string `json:"message"`
}
