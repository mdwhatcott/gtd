package date

import "time"

func YMD(ymd ...int) time.Time {
	return time.Date(ymd[0], time.Month(ymd[1]), ymd[2], 0, 0, 0, 0, time.UTC)
}