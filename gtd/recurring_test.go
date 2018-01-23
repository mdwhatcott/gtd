package gtd

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRecurringFixture(t *testing.T) {
	gunit.Run(new(RecurringFixture), t)
}

type RecurringFixture struct {
	*gunit.Fixture
	now time.Time
}

func (this *RecurringFixture) Setup() {

}

func (this *RecurringFixture) TestRecurringBothWays() {
	this.So(RecurringFromString(""), should.Equal, RecurringNever)
	this.So(RecurringFromString("asdf"), should.Equal, RecurringNever)
	this.So(RecurringFromString("monthly"), should.Equal, RecurringMonthly)
	this.So(RecurringFromString("bimonthly"), should.Equal, RecurringBimonthly)
	this.So(RecurringFromString("quarterly"), should.Equal, RecurringQuarterly)
	this.So(RecurringFromString("semiannually"), should.Equal, RecurringSemiannually)
	this.So(RecurringFromString("annually"), should.Equal, RecurringAnnually)
	this.So(RecurringFromString("biennially"), should.Equal, RecurringBiennially)

	this.So(RecurringNever.String(), should.Equal, "")
	this.So(RecurringMonthly.String(), should.Equal, "monthly")
	this.So(RecurringBimonthly.String(), should.Equal, "bimonthly")
	this.So(RecurringQuarterly.String(), should.Equal, "quarterly")
	this.So(RecurringSemiannually.String(), should.Equal, "semiannually")
	this.So(RecurringAnnually.String(), should.Equal, "annually")
	this.So(RecurringBiennially.String(), should.Equal, "biennially")
}

func (this *RecurringFixture) TestNext() {
	this.now = date(2000, 1, 1)
	this.So(RecurringNever.Next(this.now), should.Equal, date(2000, 1, 1))
	this.So(RecurringMonthly.Next(this.now), should.Equal, date(2000, 2, 1))
	this.So(RecurringBimonthly.Next(this.now), should.Equal, date(2000, 3, 1))
	this.So(RecurringQuarterly.Next(this.now), should.Equal, date(2000, 4, 1))
	this.So(RecurringSemiannually.Next(this.now), should.Equal, date(2000, 7, 1))
	this.So(RecurringAnnually.Next(this.now), should.Equal, date(2001, 1, 1))
	this.So(RecurringBiennially.Next(this.now), should.Equal, date(2002, 1, 1))
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}