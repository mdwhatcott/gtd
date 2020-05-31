package events

func TransformOutcomeTrackedV1(raw map[string]interface{}) interface{} {
	return OutcomeTrackedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		Title:     raw["title"].(string),
	}
}
func TransformOutcomeTitleUpdatedV1(raw map[string]interface{}) interface{} {
	return OutcomeTitleUpdatedV1{
		Timestamp:    ParseDate(raw["timestamp"].(string)),
		OutcomeID:    raw["outcome_id"].(string),
		UpdatedTitle: raw["updated_title"].(string),
	}
}
func TransformOutcomeExplanationUpdatedV1(raw map[string]interface{}) interface{} {
	return OutcomeExplanationUpdatedV1{
		Timestamp:          ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		UpdatedExplanation: raw["explanation"].(string),
	}
}
func TransformOutcomeDescriptionUpdatedV1(raw map[string]interface{}) interface{} {
	return OutcomeDescriptionUpdatedV1{
		Timestamp:          ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		UpdatedDescription: raw["description"].(string),
	}
}
func TransformOutcomeDeletedV1(raw map[string]interface{}) interface{} {
	return OutcomeDeletedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeFixedV1(raw map[string]interface{}) interface{} {
	return OutcomeFixedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeRealizedV1(raw map[string]interface{}) interface{} {
	return OutcomeRealizedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeAbandonedV1(raw map[string]interface{}) interface{} {
	return OutcomeAbandonedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeDeferredV1(raw map[string]interface{}) interface{} {
	return OutcomeDeferredV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformOutcomeUncertainV1(raw map[string]interface{}) interface{} {
	return OutcomeUncertainV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
	}
}
func TransformActionTrackedV1(raw map[string]interface{}) interface{} {
	return ActionTrackedV1{
		Timestamp:   ParseDate(raw["timestamp"].(string)),
		OutcomeID:   raw["outcome_id"].(string),
		ActionID:    raw["action_id"].(string),
		Description: raw["definition"].(string),
		Contexts:    TransformSlice(raw["contexts"]),
	}
}
func TransformActionsReorderedV1(raw map[string]interface{}) interface{} {
	return ActionsReorderedV1{
		Timestamp:    ParseDate(raw["timestamp"].(string)),
		OutcomeID:    raw["outcome_id"].(string),
		ReorderedIDs: TransformSlice(raw["reordered_ids"]),
	}
}
func TransformActionDescriptionUpdatedV1(raw map[string]interface{}) interface{} {
	return ActionDescriptionUpdatedV1{
		Timestamp:          ParseDate(raw["timestamp"].(string)),
		OutcomeID:          raw["outcome_id"].(string),
		ActionID:           raw["action_id"].(string),
		UpdatedDescription: raw["updated_definition"].(string),
		UpdatedContexts:    TransformSlice(raw["updated_contexts"]),
	}
}
func TransformActionStatusMarkedLatentV1(raw map[string]interface{}) interface{} {
	return ActionStatusMarkedLatentV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStatusMarkedIncompleteV1(raw map[string]interface{}) interface{} {
	return ActionStatusMarkedIncompleteV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStatusMarkedCompleteV1(raw map[string]interface{}) interface{} {
	return ActionStatusMarkedCompleteV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStrategyMarkedSequentialV1(raw map[string]interface{}) interface{} {
	return ActionStrategyMarkedSequentialV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionStrategyMarkedConcurrentV1(raw map[string]interface{}) interface{} {
	return ActionStrategyMarkedConcurrentV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
func TransformActionDeletedV1(raw map[string]interface{}) interface{} {
	return ActionDeletedV1{
		Timestamp: ParseDate(raw["timestamp"].(string)),
		OutcomeID: raw["outcome_id"].(string),
		ActionID:  raw["action_id"].(string),
	}
}
