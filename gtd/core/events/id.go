package events

func (this OutcomeTrackedV1) ID() string                 { return this.OutcomeID }
func (this OutcomeTitleUpdatedV1) ID() string            { return this.OutcomeID }
func (this OutcomeExplanationUpdatedV1) ID() string      { return this.OutcomeID }
func (this OutcomeDescriptionUpdatedV1) ID() string      { return this.OutcomeID }
func (this OutcomeDeletedV1) ID() string                 { return this.OutcomeID }
func (this OutcomeFixedV1) ID() string                   { return this.OutcomeID }
func (this OutcomeRealizedV1) ID() string                { return this.OutcomeID }
func (this OutcomeAbandonedV1) ID() string               { return this.OutcomeID }
func (this OutcomeDeferredV1) ID() string                { return this.OutcomeID }
func (this OutcomeUncertainV1) ID() string               { return this.OutcomeID }
func (this ActionTrackedV1) ID() string                  { return this.OutcomeID }
func (this ActionReorderedV1) ID() string                { return this.OutcomeID }
func (this ActionDescriptionUpdatedV1) ID() string       { return this.OutcomeID }
func (this ActionStatusMarkedLatentV1) ID() string       { return this.OutcomeID }
func (this ActionStatusMarkedIncompleteV1) ID() string   { return this.OutcomeID }
func (this ActionStatusMarkedCompleteV1) ID() string     { return this.OutcomeID }
func (this ActionStrategyMarkedSequentialV1) ID() string { return this.OutcomeID }
func (this ActionStrategyMarkedConcurrentV1) ID() string { return this.OutcomeID }
func (this ActionDeletedV1) ID() string                  { return this.OutcomeID }
