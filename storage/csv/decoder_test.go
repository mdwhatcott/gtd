package csv

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/util/fake"
)

func TestDecoderFixture(t *testing.T) {
	gunit.Run(new(DecoderFixture), t)
}

type DecoderFixture struct {
	*gunit.Fixture
	reader   *fake.Reader
	registry map[string]func([]string) interface{}
	decoder  *Decoder
}

func (this *DecoderFixture) Setup() {
	this.reader = fake.NewReader("")
	this.registry = map[string]func([]string) interface{}{
		reflect.TypeOf("").String(): func(record []string) interface{} {
			return strings.Join(record, "|")
		},
	}
	this.decoder = NewDecoder(this.reader, this.registry)
}

func (this *DecoderFixture) TestDecode() {
	this.reader.Initialize("" +
		"_timestamp_\t_id_\tstring\t_1a_\t_1b_\n" +
		"_timestamp_\t_id_\tstring\t_2a_\t_2b_\n")

	VALUE1, ERR1 := this.decoder.Decode()
	VALUE2, ERR2 := this.decoder.Decode()
	VALUE3, ERR3 := this.decoder.Decode()

	this.So(VALUE1, should.Equal, "_timestamp_|_id_|string|_1a_|_1b_")
	this.So(VALUE2, should.Equal, "_timestamp_|_id_|string|_2a_|_2b_")
	this.So(VALUE3, should.BeEmpty)

	this.So(ERR1, should.BeNil)
	this.So(ERR2, should.BeNil)
	this.So(ERR3, should.Equal, io.EOF)
}
