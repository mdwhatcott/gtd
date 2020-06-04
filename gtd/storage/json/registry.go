package json

import (
	"reflect"

	"github.com/mdwhatcott/gtd/gtd/core"
	"github.com/mdwhatcott/gtd/gtd/core/events"
)

func Registry() map[string]core.Transformer {
	return map[string]core.Transformer{
		reflect.TypeOf(events.OutcomeTrackedV1{}).String():                 TransformOutcomeTrackedV1,
		reflect.TypeOf(events.OutcomeTitleUpdatedV1{}).String():            TransformOutcomeTitleUpdatedV1,
		reflect.TypeOf(events.OutcomeExplanationUpdatedV1{}).String():      TransformOutcomeExplanationUpdatedV1,
		reflect.TypeOf(events.OutcomeDescriptionUpdatedV1{}).String():      TransformOutcomeDescriptionUpdatedV1,
		reflect.TypeOf(events.OutcomeDeletedV1{}).String():                 TransformOutcomeDeletedV1,
		reflect.TypeOf(events.OutcomeFixedV1{}).String():                   TransformOutcomeFixedV1,
		reflect.TypeOf(events.OutcomeRealizedV1{}).String():                TransformOutcomeRealizedV1,
		reflect.TypeOf(events.OutcomeAbandonedV1{}).String():               TransformOutcomeAbandonedV1,
		reflect.TypeOf(events.OutcomeDeferredV1{}).String():                TransformOutcomeDeferredV1,
		reflect.TypeOf(events.OutcomeUncertainV1{}).String():               TransformOutcomeUncertainV1,
		reflect.TypeOf(events.ActionTrackedV1{}).String():                  TransformActionTrackedV1,
		reflect.TypeOf(events.ActionsReorderedV1{}).String():               TransformActionsReorderedV1,
		reflect.TypeOf(events.ActionDescriptionUpdatedV1{}).String():       TransformActionDescriptionUpdatedV1,
		reflect.TypeOf(events.ActionStatusMarkedLatentV1{}).String():       TransformActionStatusMarkedLatentV1,
		reflect.TypeOf(events.ActionStatusMarkedIncompleteV1{}).String():   TransformActionStatusMarkedIncompleteV1,
		reflect.TypeOf(events.ActionStatusMarkedCompleteV1{}).String():     TransformActionStatusMarkedCompleteV1,
		reflect.TypeOf(events.ActionStrategyMarkedSequentialV1{}).String(): TransformActionStrategyMarkedSequentialV1,
		reflect.TypeOf(events.ActionStrategyMarkedConcurrentV1{}).String(): TransformActionStrategyMarkedConcurrentV1,
		reflect.TypeOf(events.ActionDeletedV1{}).String():                  TransformActionDeletedV1,
	}
}

func TransformOutcomeTrackedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeTrackedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		Title:     raw["title"].(string),
	}
}
func TransformOutcomeTitleUpdatedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeTitleUpdatedV1{
		Timestamp:    events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:    raw["outcome_id"].(string),
		UpdatedTitle: raw["updated_title"].(string),
	}
}
func TransformOutcomeExplanationUpdatedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeExplanationUpdatedV1{
		Timestamp:          events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		UpdatedExplanation: raw["explanation"].(string),
	}
}
func TransformOutcomeDescriptionUpdatedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeDescriptionUpdatedV1{
		Timestamp:          events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		UpdatedDescription: raw["description"].(string),
	}
}
func TransformOutcomeDeletedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeDeletedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeFixedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeFixedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeRealizedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeRealizedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeAbandonedV1(raw map[string]interface{}) interface{} {
	return events.OutcomeAbandonedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeDeferredV1(raw map[string]interface{}) interface{} {
	return events.OutcomeDeferredV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeUncertainV1(raw map[string]interface{}) interface{} {
	return events.OutcomeUncertainV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformActionTrackedV1(raw map[string]interface{}) interface{} {
	return events.ActionTrackedV1{
		Timestamp:   events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:   raw["outcome_id"].(string),
		ActionID:    raw["action_id"].(string),
		Description: raw["definition"].(string),
		Contexts:    events.TransformSlice(raw["contexts"]),
	}
}
func TransformActionsReorderedV1(raw map[string]interface{}) interface{} {
	return events.ActionsReorderedV1{
		Timestamp:    events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:    raw["outcome_id"].(string),
		ReorderedIDs: events.TransformSlice(raw["reordered_ids"]),
	}
}
func TransformActionDescriptionUpdatedV1(raw map[string]interface{}) interface{} {
	return events.ActionDescriptionUpdatedV1{
		Timestamp:          events.ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		ActionID:           raw["action_id"].(string),
		UpdatedDescription: raw["updated_definition"].(string),
		UpdatedContexts:    events.TransformSlice(raw["updated_contexts"]),
	}
}
func TransformActionStatusMarkedLatentV1(raw map[string]interface{}) interface{} {
	return events.ActionStatusMarkedLatentV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStatusMarkedIncompleteV1(raw map[string]interface{}) interface{} {
	return events.ActionStatusMarkedIncompleteV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStatusMarkedCompleteV1(raw map[string]interface{}) interface{} {
	return events.ActionStatusMarkedCompleteV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStrategyMarkedSequentialV1(raw map[string]interface{}) interface{} {
	return events.ActionStrategyMarkedSequentialV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStrategyMarkedConcurrentV1(raw map[string]interface{}) interface{} {
	return events.ActionStrategyMarkedConcurrentV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionDeletedV1(raw map[string]interface{}) interface{} {
	return events.ActionDeletedV1{
		Timestamp: events.ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
