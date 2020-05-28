package projections

import "github.com/mdwhatcott/gtd/gtd/core"

type IncompleteActionsByContext struct {
	Contexts []*Context
}

type Context struct {
	Name    string
	Actions []*ContextualAction
}

type ContextualAction struct {
	*ActionDetails
	OutcomeTitle string
}

type OutcomeDetails struct {
	Title       string
	Explanation string
	Description string
	Status      core.OutcomeStatus
	Actions     []*ActionDetails
}

type ActionDetails struct {
	ID          string
	Description string
	Contexts    []string
	Status      core.ActionStatus
	Strategy    core.ActionStrategy
}

type OutcomesListing struct {
	Fixed     []*OutcomesListingItem
	Deferred  []*OutcomesListingItem
	Uncertain []*OutcomesListingItem
	Abandoned []*OutcomesListingItem
}

type OutcomesListingItem struct {
	ID     string
	Title  string
	Status core.OutcomeStatus
}
