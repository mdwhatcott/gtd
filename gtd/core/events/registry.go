package events

import "reflect"

var types = []reflect.Type{
	reflect.TypeOf(OutcomeTrackedV1{}),
	reflect.TypeOf(OutcomeTitleUpdatedV1{}),
	reflect.TypeOf(OutcomeExplanationUpdatedV1{}),
	reflect.TypeOf(OutcomeDescriptionUpdatedV1{}),
	reflect.TypeOf(OutcomeDeletedV1{}),
	reflect.TypeOf(OutcomeFixedV1{}),
	reflect.TypeOf(OutcomeRealizedV1{}),
	reflect.TypeOf(OutcomeAbandonedV1{}),
	reflect.TypeOf(OutcomeDeferredV1{}),
	reflect.TypeOf(OutcomeUncertainV1{}),
	reflect.TypeOf(ActionTrackedV1{}),
	reflect.TypeOf(ActionReorderedV1{}),
	reflect.TypeOf(ActionDescriptionUpdatedV1{}),
	reflect.TypeOf(ActionStatusMarkedLatentV1{}),
	reflect.TypeOf(ActionStatusMarkedIncompleteV1{}),
	reflect.TypeOf(ActionStatusMarkedCompleteV1{}),
	reflect.TypeOf(ActionStrategyMarkedSequentialV1{}),
	reflect.TypeOf(ActionStrategyMarkedConcurrentV1{}),
	reflect.TypeOf(ActionDeletedV1{}),
}

func Registry() map[string]reflect.Type {
	registry := make(map[string]reflect.Type)
	for _, TYPE := range types {
		registry[TYPE.String()] = TYPE
	}
	return registry
}
