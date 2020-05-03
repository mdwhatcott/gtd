package date

import "time"

func YMD(_ymd ...int) time.Time {
	return time.Date(_ymd[0], time.Month(_ymd[1]), _ymd[2], 0, 0, 0, 0, time.UTC)
}
