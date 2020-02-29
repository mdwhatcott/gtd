package gtd

import (
	"strings"
	"time"
)

type Recurring int

const (
	RecurringNever Recurring = iota
	RecurringMonthly
	RecurringBimonthly
	RecurringQuarterly
	RecurringSemiannually
	RecurringAnnually
	RecurringBiennially

	RecurringJanuary
	RecurringFebruary
	RecurringMarch
	RecurringApril
	RecurringMay
	RecurringJune
	RecurringJuly
	RecurringAugust
	RecurringSeptember
	RecurringOctober
	RecurringNovember
	RecurringDecember
)

func FormatAllRecurringValues() string {
	var all []string
	for recurring := RecurringNever + 1; ; recurring++ {
		value := recurring.String()
		if strings.Contains(value, "invalid") {
			break
		}
		all = append(all, value)
	}
	return strings.Join(all, "|")
}

func (this Recurring) Next(now time.Time) time.Time {
	switch this {
	case RecurringJanuary:
		fallthrough
	case RecurringFebruary:
		fallthrough
	case RecurringMarch:
		fallthrough
	case RecurringApril:
		fallthrough
	case RecurringMay:
		fallthrough
	case RecurringJune:
		fallthrough
	case RecurringJuly:
		fallthrough
	case RecurringAugust:
		fallthrough
	case RecurringSeptember:
		fallthrough
	case RecurringOctober:
		fallthrough
	case RecurringNovember:
		fallthrough
	case RecurringDecember:
		return nextYear(now, recurringMonths[this])
	case RecurringBiennially:
		return now.AddDate(2, 0, 0)
	case RecurringAnnually:
		return now.AddDate(1, 0, 0)
	case RecurringSemiannually:
		return now.AddDate(0, 6, 0)
	case RecurringQuarterly:
		return now.AddDate(0, 3, 0)
	case RecurringBimonthly:
		return now.AddDate(0, 2, 0)
	case RecurringMonthly:
		return now.AddDate(0, 1, 0)
	default:
		return now.AddDate(0, 0, 0)
	}
}

func nextYear(now time.Time, target time.Month) time.Time {
	if now.Month() == target {
		now = now.AddDate(1, 0, 0)
	}
	for now.Month() != target {
		now = now.AddDate(0, 1, 0)
	}
	return now
}

func (this Recurring) String() string {
	if value, found := recurringMap[this]; found {
		return value.(string)
	} else {
		return "invalid (Recurring)"
	}
}

func RecurringFromString(raw string) Recurring {
	if value, found := recurringMap[raw]; found {
		return value.(Recurring)
	} else {
		return RecurringNever
	}
}

var recurringMap = map[interface{}]interface{}{
	"monthly":               RecurringMonthly,
	"bimonthly":             RecurringBimonthly,
	"quarterly":             RecurringQuarterly,
	"semiannually":          RecurringSemiannually,
	"annually":              RecurringAnnually,
	"biennially":            RecurringBiennially,
	time.January.String():   RecurringJanuary,
	time.February.String():  RecurringFebruary,
	time.March.String():     RecurringMarch,
	time.April.String():     RecurringApril,
	time.May.String():       RecurringMay,
	time.June.String():      RecurringJune,
	time.July.String():      RecurringJuly,
	time.August.String():    RecurringAugust,
	time.September.String(): RecurringSeptember,
	time.October.String():   RecurringOctober,
	time.November.String():  RecurringNovember,
	time.December.String():  RecurringDecember,
	RecurringNever:          "",
	RecurringMonthly:        "monthly",
	RecurringBimonthly:      "bimonthly",
	RecurringQuarterly:      "quarterly",
	RecurringSemiannually:   "semiannually",
	RecurringAnnually:       "annually",
	RecurringBiennially:     "biennially",
	RecurringJanuary:        time.January.String(),
	RecurringFebruary:       time.February.String(),
	RecurringMarch:          time.March.String(),
	RecurringApril:          time.April.String(),
	RecurringMay:            time.May.String(),
	RecurringJune:           time.June.String(),
	RecurringJuly:           time.July.String(),
	RecurringAugust:         time.August.String(),
	RecurringSeptember:      time.September.String(),
	RecurringOctober:        time.October.String(),
	RecurringNovember:       time.November.String(),
	RecurringDecember:       time.December.String(),
}

var recurringMonths = map[Recurring]time.Month{
	RecurringJanuary:   time.January,
	RecurringFebruary:  time.February,
	RecurringMarch:     time.March,
	RecurringApril:     time.April,
	RecurringMay:       time.May,
	RecurringJune:      time.June,
	RecurringJuly:      time.July,
	RecurringAugust:    time.August,
	RecurringSeptember: time.September,
	RecurringOctober:   time.October,
	RecurringNovember:  time.November,
	RecurringDecember:  time.December,
}
