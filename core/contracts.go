package core

import (
	"context"
	"errors"
	"time"
)

type Handler interface {
	Handle(context.Context, ...interface{})
}

type Projector interface {
	Apply(events chan interface{})
	Projection() interface{}
}

type IDFunc func() string

var (
	ErrOutcomeNotFound  = errors.New("outcome not found")
	ErrActionNotFound   = errors.New("action not found")
	ErrOutcomeUnchanged = errors.New("outcome unchanged")
)

type OutcomeStatus string

const (
	OutcomeStatusFixed     OutcomeStatus = "FIXED"
	OutcomeStatusRealized  OutcomeStatus = "REALIZED"
	OutcomeStatusAbandoned OutcomeStatus = "ABANDONED"
	OutcomeStatusDeferred  OutcomeStatus = "DEFERRED"
	OutcomeStatusUncertain OutcomeStatus = "UNCERTAIN"
)

type ActionStatus string

const (
	ActionStatusIncomplete ActionStatus = "INCOMPLETE"
	ActionStatusLatent     ActionStatus = "LATENT"
	ActionStatusComplete   ActionStatus = "COMPLETE"
)

type ActionStrategy string

const (
	ActionStrategyConcurrent ActionStrategy = "CONCURRENT"
	ActionStrategySequential ActionStrategy = "SEQUENTIAL"
)

type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

type Clock func() time.Time

func Now() time.Time {
	return time.Now().UTC()
}
