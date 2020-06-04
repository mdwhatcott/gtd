package events

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/util/date"
)

type Time = time.Time

var (
	ParseDateRFC3339     = date.ParseRFC3339
	ParseDateRFC3339Nano = date.ParseRFC3339Nano
)

func TransformSlice(raw interface{}) (transformed []string) {
	if raw == nil {
		return transformed
	}
	for _, item := range raw.([]interface{}) {
		transformed = append(transformed, item.(string))
	}
	return transformed
}
