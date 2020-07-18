package csv

import (
	"reflect"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"

	"github.com/mdwhatcott/gtd/v3/util/fake"
)

func TestEncoderFixture(t *testing.T) {
	gunit.Run(new(EncoderFixture), t)
}

type EncoderFixture struct {
	*gunit.Fixture

	writer   *fake.Writer
	encoder  *Encoder
	registry map[string]func(interface{}) []string
}

func (this *EncoderFixture) Setup() {
	this.writer = fake.NewWriter(nil, nil)
	this.registry = map[string]func(interface{}) []string{
		reflect.TypeOf("").String(): func(i interface{}) []string { return strings.Fields(i.(string)) },
	}
	this.encoder = NewEncoder(this.writer, this.registry)
}

func (this *EncoderFixture) TestEncode() {
	err1 := this.encoder.Encode("this is a test")
	err2 := this.encoder.Encode("this is a 2nd test")

	this.So(err1, should.BeNil)
	this.So(err2, should.BeNil)
	this.So(this.writer.Content(), should.Equal, ""+
		"this\tis\ta\ttest\n"+
		"this\tis\ta\t2nd\ttest\n")
}
