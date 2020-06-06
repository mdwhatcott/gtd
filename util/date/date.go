package date

import "time"

func ParseRFC3339Nano(raw string) time.Time {
	parsed, _ := time.Parse(time.RFC3339Nano, raw)
	return parsed
}
