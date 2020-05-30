package events

import (
	"time"

	"github.com/mdwhatcott/gtd/gtd/util/date"
)

type Time = time.Time

var ParseDate = date.Parse

func TransformSlice(raw []interface{}) (transformed []string) {
	for _, item := range raw {
		transformed = append(transformed, item.(string))
	}
	return transformed
}
