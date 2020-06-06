package domain

import "github.com/mdwhatcott/gtd/core"

type Action struct {
	Description string
	Status      core.ActionStatus
	Strategy    core.ActionStrategy
}
