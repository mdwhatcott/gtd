package domain

import "github.com/mdwhatcott/gtd/v3/core"

type Action struct {
	Description string
	Status      core.ActionStatus
	Strategy    core.ActionStrategy
}
