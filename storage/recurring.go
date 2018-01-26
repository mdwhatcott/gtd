package storage

import "time"


type Recurring int

const (
	RecurringNever        Recurring = iota
	RecurringMonthly
	RecurringBimonthly
	RecurringQuarterly
	RecurringSemiannually
	RecurringAnnually
	RecurringBiennially
)

func (this Recurring) Next(now time.Time) time.Time {
	switch this {
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
	"monthly":             RecurringMonthly,
	"bimonthly":           RecurringBimonthly,
	"quarterly":           RecurringQuarterly,
	"semiannually":        RecurringSemiannually,
	"annually":            RecurringAnnually,
	"biennially":          RecurringBiennially,
	RecurringNever:        "",
	RecurringMonthly:      "monthly",
	RecurringBimonthly:    "bimonthly",
	RecurringQuarterly:    "quarterly",
	RecurringSemiannually: "semiannually",
	RecurringAnnually:     "annually",
	RecurringBiennially:   "biennially",
}

