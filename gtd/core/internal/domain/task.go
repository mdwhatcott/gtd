package domain

import "github.com/smartystreets/joyride/v2"

type Task struct {
	*joyride.Base
}

func NewTask(string) *Task {
	return &Task{Base: joyride.New()}
}
