package outcomes

import "github.com/mdwhatcott/gtd/gtd/core"

type Action struct {
	ID          string
	Description string
	Status      core.ActionStatus
}
