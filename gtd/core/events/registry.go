package events

import (
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/core"
)

func Registry() map[string]core.Transformer {
	return map[string]core.Transformer{
		reflect.TypeOf(OutcomeTrackedV1{}).String():                 TransformOutcomeTrackedV1,
		reflect.TypeOf(OutcomeTitleUpdatedV1{}).String():            TransformOutcomeTitleUpdatedV1,
		reflect.TypeOf(OutcomeExplanationUpdatedV1{}).String():      TransformOutcomeExplanationUpdatedV1,
		reflect.TypeOf(OutcomeDescriptionUpdatedV1{}).String():      TransformOutcomeDescriptionUpdatedV1,
		reflect.TypeOf(OutcomeDeletedV1{}).String():                 TransformOutcomeDeletedV1,
		reflect.TypeOf(OutcomeFixedV1{}).String():                   TransformOutcomeFixedV1,
		reflect.TypeOf(OutcomeRealizedV1{}).String():                TransformOutcomeRealizedV1,
		reflect.TypeOf(OutcomeAbandonedV1{}).String():               TransformOutcomeAbandonedV1,
		reflect.TypeOf(OutcomeDeferredV1{}).String():                TransformOutcomeDeferredV1,
		reflect.TypeOf(OutcomeUncertainV1{}).String():               TransformOutcomeUncertainV1,
		reflect.TypeOf(ActionTrackedV1{}).String():                  TransformActionTrackedV1,
		reflect.TypeOf(ActionsReorderedV1{}).String():               TransformActionsReorderedV1,
		reflect.TypeOf(ActionDescriptionUpdatedV1{}).String():       TransformActionDescriptionUpdatedV1,
		reflect.TypeOf(ActionStatusMarkedLatentV1{}).String():       TransformActionStatusMarkedLatentV1,
		reflect.TypeOf(ActionStatusMarkedIncompleteV1{}).String():   TransformActionStatusMarkedIncompleteV1,
		reflect.TypeOf(ActionStatusMarkedCompleteV1{}).String():     TransformActionStatusMarkedCompleteV1,
		reflect.TypeOf(ActionStrategyMarkedSequentialV1{}).String(): TransformActionStrategyMarkedSequentialV1,
		reflect.TypeOf(ActionStrategyMarkedConcurrentV1{}).String(): TransformActionStrategyMarkedConcurrentV1,
		reflect.TypeOf(ActionDeletedV1{}).String():                  TransformActionDeletedV1,
	}
}
