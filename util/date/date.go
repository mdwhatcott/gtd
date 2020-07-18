package date

import "time"

func ParseRFC3339Nano(raw string) time.Time {
	PARSED, _ := time.Parse(time.RFC3339Nano, raw)
	return PARSED
}

func Today() string {
	return time.Now().Format("2006-01-02")
}
