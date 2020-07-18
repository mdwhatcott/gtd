package csv

import (
	"strings"
	"time"

	"github.com/mdwhatcott/gtd/v3/core/events"
	"github.com/mdwhatcott/gtd/v3/util/date"
)

func EncoderRegistry() map[string]func(interface{}) []string {
	return map[string]func(interface{}) []string{
		"events.OutcomeTrackedV1":                 EncodeOutcomeTrackedV1,
		"events.OutcomeTitleUpdatedV1":            EncodeOutcomeTitleUpdatedV1,
		"events.OutcomeExplanationUpdatedV1":      EncodeOutcomeExplanationUpdatedV1,
		"events.OutcomeDescriptionUpdatedV1":      EncodeOutcomeDescriptionUpdatedV1,
		"events.OutcomeDeletedV1":                 EncodeOutcomeDeletedV1,
		"events.OutcomeFixedV1":                   EncodeOutcomeFixedV1,
		"events.OutcomeRealizedV1":                EncodeOutcomeRealizedV1,
		"events.OutcomeAbandonedV1":               EncodeOutcomeAbandonedV1,
		"events.OutcomeDeferredV1":                EncodeOutcomeDeferredV1,
		"events.OutcomeUncertainV1":               EncodeOutcomeUncertainV1,
		"events.ActionTrackedV1":                  EncodeActionTrackedV1,
		"events.ActionsReorderedV1":               EncodeActionsReorderedV1,
		"events.ActionDescriptionUpdatedV1":       EncodeActionDescriptionUpdatedV1,
		"events.ActionStatusMarkedLatentV1":       EncodeActionStatusMarkedLatentV1,
		"events.ActionStatusMarkedIncompleteV1":   EncodeActionStatusMarkedIncompleteV1,
		"events.ActionStatusMarkedCompleteV1":     EncodeActionStatusMarkedCompleteV1,
		"events.ActionStrategyMarkedSequentialV1": EncodeActionStrategyMarkedSequentialV1,
		"events.ActionStrategyMarkedConcurrentV1": EncodeActionStrategyMarkedConcurrentV1,
		"events.ActionDeletedV1":                  EncodeActionDeletedV1,
	}
}
func DecoderRegistry() map[string]func([]string) interface{} {
	return map[string]func([]string) interface{}{
		"events.OutcomeTrackedV1":                 DecodeOutcomeTrackedV1,
		"events.OutcomeTitleUpdatedV1":            DecodeOutcomeTitleUpdatedV1,
		"events.OutcomeExplanationUpdatedV1":      DecodeOutcomeExplanationUpdatedV1,
		"events.OutcomeDescriptionUpdatedV1":      DecodeOutcomeDescriptionUpdatedV1,
		"events.OutcomeDeletedV1":                 DecodeOutcomeDeletedV1,
		"events.OutcomeFixedV1":                   DecodeOutcomeFixedV1,
		"events.OutcomeRealizedV1":                DecodeOutcomeRealizedV1,
		"events.OutcomeAbandonedV1":               DecodeOutcomeAbandonedV1,
		"events.OutcomeDeferredV1":                DecodeOutcomeDeferredV1,
		"events.OutcomeUncertainV1":               DecodeOutcomeUncertainV1,
		"events.ActionTrackedV1":                  DecodeActionTrackedV1,
		"events.ActionsReorderedV1":               DecodeActionsReorderedV1,
		"events.ActionDescriptionUpdatedV1":       DecodeActionDescriptionUpdatedV1,
		"events.ActionStatusMarkedLatentV1":       DecodeActionStatusMarkedLatentV1,
		"events.ActionStatusMarkedIncompleteV1":   DecodeActionStatusMarkedIncompleteV1,
		"events.ActionStatusMarkedCompleteV1":     DecodeActionStatusMarkedCompleteV1,
		"events.ActionStrategyMarkedSequentialV1": DecodeActionStrategyMarkedSequentialV1,
		"events.ActionStrategyMarkedConcurrentV1": DecodeActionStrategyMarkedConcurrentV1,
		"events.ActionDeletedV1":                  DecodeActionDeletedV1,
	}
}

func EncodeOutcomeTrackedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeTrackedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeTrackedV1",
		EVENT.Title,
	}
}
func EncodeOutcomeTitleUpdatedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeTitleUpdatedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeTitleUpdatedV1",
		EVENT.UpdatedTitle,
	}
}
func EncodeOutcomeExplanationUpdatedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeExplanationUpdatedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeExplanationUpdatedV1",
		EVENT.UpdatedExplanation,
	}
}
func EncodeOutcomeDescriptionUpdatedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeDescriptionUpdatedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeDescriptionUpdatedV1",
		strings.ReplaceAll(EVENT.UpdatedDescription, "\n", "Ω"),
	}
}
func EncodeOutcomeDeletedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeDeletedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeDeletedV1",
	}
}
func EncodeOutcomeFixedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeFixedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeFixedV1",
	}
}
func EncodeOutcomeRealizedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeRealizedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeRealizedV1",
	}
}
func EncodeOutcomeAbandonedV1(v interface{}) []string {
	EVENT := v.(events.OutcomeAbandonedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeAbandonedV1",
	}
}
func EncodeOutcomeDeferredV1(v interface{}) []string {
	EVENT := v.(events.OutcomeDeferredV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeDeferredV1",
	}
}
func EncodeOutcomeUncertainV1(v interface{}) []string {
	EVENT := v.(events.OutcomeUncertainV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.OutcomeUncertainV1",
	}
}
func EncodeActionTrackedV1(v interface{}) []string {
	EVENT := v.(events.ActionTrackedV1)
	var contexts string
	if len(EVENT.Contexts) > 0 {
		contexts = strings.Join(EVENT.Contexts, "|")
	}
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionTrackedV1",
		EVENT.ActionID,
		EVENT.Description,
		contexts,
	}
}
func EncodeActionsReorderedV1(v interface{}) []string {
	EVENT := v.(events.ActionsReorderedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionsReorderedV1",
		strings.Join(EVENT.ReorderedIDs, "|"),
	}
}
func EncodeActionDescriptionUpdatedV1(v interface{}) []string {
	EVENT := v.(events.ActionDescriptionUpdatedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionDescriptionUpdatedV1",
		EVENT.ActionID,
		EVENT.UpdatedDescription,
		strings.Join(EVENT.UpdatedContexts, "|"),
	}
}
func EncodeActionStatusMarkedLatentV1(v interface{}) []string {
	EVENT := v.(events.ActionStatusMarkedLatentV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionStatusMarkedLatentV1",
		EVENT.ActionID,
	}
}
func EncodeActionStatusMarkedIncompleteV1(v interface{}) []string {
	EVENT := v.(events.ActionStatusMarkedIncompleteV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionStatusMarkedIncompleteV1",
		EVENT.ActionID,
	}
}
func EncodeActionStatusMarkedCompleteV1(v interface{}) []string {
	EVENT := v.(events.ActionStatusMarkedCompleteV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionStatusMarkedCompleteV1",
		EVENT.ActionID,
	}
}
func EncodeActionStrategyMarkedSequentialV1(v interface{}) []string {
	EVENT := v.(events.ActionStrategyMarkedSequentialV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionStrategyMarkedSequentialV1",
		EVENT.ActionID,
	}
}
func EncodeActionStrategyMarkedConcurrentV1(v interface{}) []string {
	EVENT := v.(events.ActionStrategyMarkedConcurrentV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionStrategyMarkedConcurrentV1",
		EVENT.ActionID,
	}
}
func EncodeActionDeletedV1(v interface{}) []string {
	EVENT := v.(events.ActionDeletedV1)
	return []string{
		EVENT.Timestamp.Format(time.RFC3339Nano),
		EVENT.OutcomeID,
		"events.ActionDeletedV1",
		EVENT.ActionID,
	}
}

func DecodeOutcomeTrackedV1(record []string) interface{} {
	return events.OutcomeTrackedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		Title:     record[3],
	}
}
func DecodeOutcomeTitleUpdatedV1(record []string) interface{} {
	return events.OutcomeTitleUpdatedV1{
		Timestamp:    date.ParseRFC3339Nano(record[0]),
		OutcomeID:    record[1],
		UpdatedTitle: record[3],
	}
}
func DecodeOutcomeExplanationUpdatedV1(record []string) interface{} {
	return events.OutcomeExplanationUpdatedV1{
		Timestamp:          date.ParseRFC3339Nano(record[0]),
		OutcomeID:          record[1],
		UpdatedExplanation: record[3],
	}
}
func DecodeOutcomeDescriptionUpdatedV1(record []string) interface{} {
	return events.OutcomeDescriptionUpdatedV1{
		Timestamp:          date.ParseRFC3339Nano(record[0]),
		OutcomeID:          record[1],
		UpdatedDescription: strings.ReplaceAll(record[3], "Ω", "\n"),
	}
}
func DecodeOutcomeDeletedV1(record []string) interface{} {
	return events.OutcomeDeletedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeOutcomeFixedV1(record []string) interface{} {
	return events.OutcomeFixedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeOutcomeRealizedV1(record []string) interface{} {
	return events.OutcomeRealizedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeOutcomeAbandonedV1(record []string) interface{} {
	return events.OutcomeAbandonedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeOutcomeDeferredV1(record []string) interface{} {
	return events.OutcomeDeferredV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeOutcomeUncertainV1(record []string) interface{} {
	return events.OutcomeUncertainV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
	}
}
func DecodeActionTrackedV1(record []string) interface{} {
	values := strings.Split(record[5], "|")
	if len(values) == 1 && values[0] == "" {
		values = nil
	}
	return events.ActionTrackedV1{
		Timestamp:   date.ParseRFC3339Nano(record[0]),
		OutcomeID:   record[1],
		ActionID:    record[3],
		Description: record[4],
		Contexts:    decodeSlice(record[5]),
	}
}
func DecodeActionsReorderedV1(record []string) interface{} {
	return events.ActionsReorderedV1{
		Timestamp:    date.ParseRFC3339Nano(record[0]),
		OutcomeID:    record[1],
		ReorderedIDs: decodeSlice(record[3]),
	}
}
func DecodeActionDescriptionUpdatedV1(record []string) interface{} {
	return events.ActionDescriptionUpdatedV1{
		Timestamp:          date.ParseRFC3339Nano(record[0]),
		OutcomeID:          record[1],
		ActionID:           record[3],
		UpdatedDescription: record[4],
		UpdatedContexts:    decodeSlice(record[5]),
	}
}
func DecodeActionStatusMarkedLatentV1(record []string) interface{} {
	return events.ActionStatusMarkedLatentV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}
func DecodeActionStatusMarkedIncompleteV1(record []string) interface{} {
	return events.ActionStatusMarkedIncompleteV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}
func DecodeActionStatusMarkedCompleteV1(record []string) interface{} {
	return events.ActionStatusMarkedCompleteV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}
func DecodeActionStrategyMarkedSequentialV1(record []string) interface{} {
	return events.ActionStrategyMarkedSequentialV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}
func DecodeActionStrategyMarkedConcurrentV1(record []string) interface{} {
	return events.ActionStrategyMarkedConcurrentV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}
func DecodeActionDeletedV1(record []string) interface{} {
	return events.ActionDeletedV1{
		Timestamp: date.ParseRFC3339Nano(record[0]),
		OutcomeID: record[1],
		ActionID:  record[3],
	}
}

func decodeSlice(field string) []string {
	VALUES := strings.Split(field, "|")
	if len(VALUES) == 1 && VALUES[0] == "" {
		VALUES = nil
	}
	return VALUES
}
