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

func (this *RecurringFixture) TestRecurringBothWays() {
	this.So(RecurringFromString(""), should.Equal, RecurringNever)
	this.So(RecurringFromString("asdf"), should.Equal, RecurringNever)
	this.So(RecurringFromString("monthly"), should.Equal, RecurringMonthly)
	this.So(RecurringFromString("bimonthly"), should.Equal, RecurringBimonthly)
	this.So(RecurringFromString("quarterly"), should.Equal, RecurringQuarterly)
	this.So(RecurringFromString("semiannually"), should.Equal, RecurringSemiannually)
	this.So(RecurringFromString("annually"), should.Equal, RecurringAnnually)
	this.So(RecurringFromString("biennially"), should.Equal, RecurringBiennially)
	this.So(RecurringFromString("January"), should.Equal, RecurringJanuary)
	this.So(RecurringFromString("February"), should.Equal, RecurringFebruary)
	this.So(RecurringFromString("March"), should.Equal, RecurringMarch)
	this.So(RecurringFromString("April"), should.Equal, RecurringApril)
	this.So(RecurringFromString("May"), should.Equal, RecurringMay)
	this.So(RecurringFromString("June"), should.Equal, RecurringJune)
	this.So(RecurringFromString("July"), should.Equal, RecurringJuly)
	this.So(RecurringFromString("August"), should.Equal, RecurringAugust)
	this.So(RecurringFromString("September"), should.Equal, RecurringSeptember)
	this.So(RecurringFromString("October"), should.Equal, RecurringOctober)
	this.So(RecurringFromString("November"), should.Equal, RecurringNovember)
	this.So(RecurringFromString("December"), should.Equal, RecurringDecember)

	this.So(RecurringNever.String(), should.Equal, "")
	this.So(RecurringMonthly.String(), should.Equal, "monthly")
	this.So(RecurringBimonthly.String(), should.Equal, "bimonthly")
	this.So(RecurringQuarterly.String(), should.Equal, "quarterly")
	this.So(RecurringSemiannually.String(), should.Equal, "semiannually")
	this.So(RecurringAnnually.String(), should.Equal, "annually")
	this.So(RecurringBiennially.String(), should.Equal, "biennially")
	this.So(RecurringJanuary.String(), should.Equal, "January")
	this.So(RecurringFebruary.String(), should.Equal, "February")
	this.So(RecurringMarch.String(), should.Equal, "March")
	this.So(RecurringApril.String(), should.Equal, "April")
	this.So(RecurringMay.String(), should.Equal, "May")
	this.So(RecurringJune.String(), should.Equal, "June")
	this.So(RecurringJuly.String(), should.Equal, "July")
	this.So(RecurringAugust.String(), should.Equal, "August")
	this.So(RecurringSeptember.String(), should.Equal, "September")
	this.So(RecurringOctober.String(), should.Equal, "October")
	this.So(RecurringNovember.String(), should.Equal, "November")
	this.So(RecurringDecember.String(), should.Equal, "December")
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

func (this *RecurringFixture) TestNext_FixedYearly() {
	this.So(RecurringJanuary.Next(date(2000, 1, 1)), should.Equal, date(2001, 1, 1))
	this.So(RecurringJanuary.Next(date(2000, 2, 1)), should.Equal, date(2001, 1, 1))

	this.So(RecurringFebruary.Next(date(2000, 1, 1)), should.Equal, date(2000, 2, 1))
	this.So(RecurringFebruary.Next(date(2000, 2, 1)), should.Equal, date(2001, 2, 1))

	this.So(RecurringMarch.Next(date(2000, 1, 1)), should.Equal, date(2000, 3, 1))
	this.So(RecurringMarch.Next(date(2000, 3, 1)), should.Equal, date(2001, 3, 1))

	this.So(RecurringApril.Next(date(2000, 1, 1)), should.Equal, date(2000, 4, 1))
	this.So(RecurringApril.Next(date(2000, 4, 1)), should.Equal, date(2001, 4, 1))

	this.So(RecurringMay.Next(date(2000, 1, 1)), should.Equal, date(2000, 5, 1))
	this.So(RecurringMay.Next(date(2000, 5, 1)), should.Equal, date(2001, 5, 1))

	this.So(RecurringJune.Next(date(2000, 1, 1)), should.Equal, date(2000, 6, 1))
	this.So(RecurringJune.Next(date(2000, 6, 1)), should.Equal, date(2001, 6, 1))

	this.So(RecurringJuly.Next(date(2000, 1, 1)), should.Equal, date(2000, 7, 1))
	this.So(RecurringJuly.Next(date(2000, 7, 1)), should.Equal, date(2001, 7, 1))

	this.So(RecurringAugust.Next(date(2000, 1, 1)), should.Equal, date(2000, 8, 1))
	this.So(RecurringAugust.Next(date(2000, 8, 1)), should.Equal, date(2001, 8, 1))

	this.So(RecurringSeptember.Next(date(2000, 1, 1)), should.Equal, date(2000, 9, 1))
	this.So(RecurringSeptember.Next(date(2000, 9, 1)), should.Equal, date(2001, 9, 1))

	this.So(RecurringOctober.Next(date(2000, 1, 1)), should.Equal, date(2000, 10, 1))
	this.So(RecurringOctober.Next(date(2000, 10, 1)), should.Equal, date(2001, 10, 1))

	this.So(RecurringNovember.Next(date(2000, 1, 1)), should.Equal, date(2000, 11, 1))
	this.So(RecurringNovember.Next(date(2000, 11, 1)), should.Equal, date(2001, 11, 1))

	this.So(RecurringDecember.Next(date(2000, 1, 1)), should.Equal, date(2000, 12, 1))
	this.So(RecurringDecember.Next(date(2000, 12, 1)), should.Equal, date(2001, 12, 1))
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
