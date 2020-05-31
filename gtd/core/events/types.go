package events

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/util/date"
)

type Time = time.Time

var ParseDate = date.Parse

func TransformSlice(raw interface{}) (transformed []string) {
	if raw == nil {
		return transformed
	}
	for _, item := range raw.([]interface{}) {
		transformed = append(transformed, item.(string))
	}
	return transformed
}
