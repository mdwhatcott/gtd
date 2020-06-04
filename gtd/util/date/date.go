package date

import "time"

func YMD(_ymd ...int) time.Time {
	return time.Date(_ymd[0], time.Month(_ymd[1]), _ymd[2], 0, 0, 0, 0, time.UTC)
}

func ParseRFC3339(raw string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, raw)
	return parsed
}

func ParseRFC3339Nano(raw string) time.Time {
	parsed, _ := time.Parse(time.RFC3339Nano, raw)
	return parsed
}
