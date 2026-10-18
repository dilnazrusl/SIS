package task

import (
	"fmt"
	"slices"
	"strings"
)

type Manager struct {
	processor Processor
	tasks     map[int]*Task
	nextID    int
}

func NewManager(p Processor) *Manager {
	if p == nil {
		p = TextAnalyzer{}
	}
	return &Manager{
		processor: p,
		tasks:     make(map[int]*Task),
		nextID:    1,
	}
}

func (m *Manager) Create(input string) (int, error) {
	if err := ValidateInput(input); err != nil {
		return 0, err
	}

	t := &Task{
		ID:    m.nextID,
		Input: strings.TrimSpace(input),
		State: StatePending,
	}
	m.tasks[t.ID] = t
	m.nextID++

	return t.ID, nil
}

func (m *Manager) find(id int) (*Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
	}
	return t, nil
}

func (m *Manager) Get(id int) (Task, error) {
	t, err := m.find(id)
	if err != nil {
		return Task{}, err
	}
	return *t, nil
}

func (m *Manager) List() []Task {
	list := make([]Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		list = append(list, *t)
	}
	slices.SortFunc(list, func(a, b Task) int { return a.ID - b.ID })
	return list
}

func (m *Manager) Process(id int) error {
	t, err := m.find(id)
	if err != nil {
		return err
	}

	if err := t.transition(StateProcessing); err != nil {
		return fmt.Errorf("task %d: %w", id, err)
	}

	result, err := m.processor.Process(t.Input)
	if err != nil {
		if terr := t.transition(StateFailed); terr != nil {
			return terr
		}
		return fmt.Errorf("process task %d: %w", id, err)
	}

	t.Result = result
	return t.transition(StateCompleted)
}
