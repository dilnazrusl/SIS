package task

import (
	"fmt"
	"slices"
)

type State string

const (
	StatePending    State = "Pending"
	StateProcessing State = "Processing"
	StateCompleted  State = "Completed"
	StateFailed     State = "Failed"
)

var allowedTransitions = map[State][]State{
	StatePending:    {StateProcessing},
	StateProcessing: {StateCompleted, StateFailed},
}

type Result struct {
	Words int
	Runes int
	Bytes int
}

type Task struct {
	ID     int
	Input  string
	State  State
	Result Result
}

func (t *Task) transition(to State) error {
	if !slices.Contains(allowedTransitions[t.State], to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, t.State, to)
	}
	t.State = to
	return nil
}
