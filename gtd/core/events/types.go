package events

import "time"

type Time = time.Time

type ActionStatus int

const (
	ActionLatent ActionStatus = iota
	ActionPending
	ActionComplete
)

type ActionStrategy int

const (
	ActionConcurrent ActionStrategy = iota
	ActionSequential
)
