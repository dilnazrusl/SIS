package task

import (
	"errors"
	"testing"
)

var errBoom = errors.New("boom")

type failingProcessor struct{}

func (failingProcessor) Process(string) (Result, error) {
	return Result{}, errBoom
}

func TestCreateTask(t *testing.T) {
	m := NewManager(TextAnalyzer{})

	id, err := m.Create("Hello Go world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("expected first ID 1, got %d", id)
	}

	got, err := m.Get(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.State != StatePending {
		t.Errorf("expected state %s, got %s", StatePending, got.State)
	}
	if got.Input != "Hello Go world" {
		t.Errorf("unexpected input %q", got.Input)
	}
	if got.Result != (Result{}) {
		t.Errorf("expected empty result before processing, got %+v", got.Result)
	}
}

func TestCreateTaskTrimsInput(t *testing.T) {
	m := NewManager(TextAnalyzer{})

	id, err := m.Create("  hello  world  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := m.Get(id)
	if got.Input != "hello  world" {
		t.Errorf("expected trimmed input, got %q", got.Input)
	}
}

func TestCreateTaskInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"spaces", "   "},
		{"tabs and newlines", "\t\n  \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(TextAnalyzer{})

			id, err := m.Create(tt.input)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("expected ErrInvalidInput, got %v", err)
			}
			if id != 0 {
				t.Errorf("expected ID 0 on error, got %d", id)
			}
			if n := len(m.List()); n != 0 {
				t.Errorf("invalid task must not be stored, got %d tasks", n)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	m := NewManager(TextAnalyzer{})
	id, _ := m.Create("Hello Go world")

	got, err := m.Get(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected ID %d, got %d", id, got.ID)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	m := NewManager(TextAnalyzer{})

	for _, id := range []int{0, -1, 999} {
		_, err := m.Get(id)
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("id %d: expected ErrTaskNotFound, got %v", id, err)
		}
	}
}

func TestGetReturnsCopy(t *testing.T) {
	m := NewManager(TextAnalyzer{})
	id, _ := m.Create("Hello Go world")

	got, _ := m.Get(id)
	got.State = StateCompleted
	got.Input = "changed"

	again, _ := m.Get(id)
	if again.State != StatePending || again.Input != "Hello Go world" {
		t.Errorf("stored task was modified through a returned copy: %+v", again)
	}
}

func TestProcessTask(t *testing.T) {
	m := NewManager(TextAnalyzer{})
	id, _ := m.Create("Hello Go world")

	if err := m.Process(id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := m.Get(id)
	if got.State != StateCompleted {
		t.Errorf("expected state %s, got %s", StateCompleted, got.State)
	}
	want := Result{Words: 3, Runes: 14, Bytes: 14}
	if got.Result != want {
		t.Errorf("expected result %+v, got %+v", want, got.Result)
	}
}

func TestProcessTaskResults(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Result
	}{
		{"single word", "Go", Result{Words: 1, Runes: 2, Bytes: 2}},
		{"three words", "Hello Go world", Result{Words: 3, Runes: 14, Bytes: 14}},
		{"cyrillic", "Привет мир", Result{Words: 2, Runes: 10, Bytes: 19}},
		{"extra spaces between words", "a   b    c", Result{Words: 3, Runes: 10, Bytes: 10}},
		{"surrounding spaces trimmed", "  Go  ", Result{Words: 1, Runes: 2, Bytes: 2}},
		{"multi-byte symbols", "世界", Result{Words: 1, Runes: 2, Bytes: 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(TextAnalyzer{})
			id, err := m.Create(tt.input)
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			if err := m.Process(id); err != nil {
				t.Fatalf("process: %v", err)
			}
			got, _ := m.Get(id)
			if got.Result != tt.want {
				t.Errorf("expected %+v, got %+v", tt.want, got.Result)
			}
			if got.State != StateCompleted {
				t.Errorf("expected state %s, got %s", StateCompleted, got.State)
			}
		})
	}
}

func TestProcessNotFound(t *testing.T) {
	m := NewManager(TextAnalyzer{})

	err := m.Process(42)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestProcessAlreadyCompleted(t *testing.T) {
	m := NewManager(TextAnalyzer{})
	id, _ := m.Create("Hello Go world")
	if err := m.Process(id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	before, _ := m.Get(id)

	err := m.Process(id)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}

	after, _ := m.Get(id)
	if after != before {
		t.Errorf("rejected operation changed the task: before %+v, after %+v", before, after)
	}
}

func TestProcessFailure(t *testing.T) {
	m := NewManager(failingProcessor{})
	id, _ := m.Create("Hello Go world")

	err := m.Process(id)
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped errBoom, got %v", err)
	}

	got, _ := m.Get(id)
	if got.State != StateFailed {
		t.Errorf("expected state %s, got %s", StateFailed, got.State)
	}
	if got.Result != (Result{}) {
		t.Errorf("failed task must have empty result, got %+v", got.Result)
	}

	if err := m.Process(id); !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition for failed task, got %v", err)
	}
}

func TestMultipleIndependentTasks(t *testing.T) {
	m := NewManager(TextAnalyzer{})

	inputs := []string{"Go", "Hello Go world", "Привет мир", "one two three four five"}
	ids := make([]int, 0, len(inputs))
	seen := make(map[int]bool)

	for _, in := range inputs {
		id, err := m.Create(in)
		if err != nil {
			t.Fatalf("create %q: %v", in, err)
		}
		if seen[id] {
			t.Fatalf("duplicate ID %d", id)
		}
		seen[id] = true
		ids = append(ids, id)
	}

	if err := m.Process(ids[1]); err != nil {
		t.Fatalf("process: %v", err)
	}

	for i, id := range ids {
		got, _ := m.Get(id)
		if got.Input != inputs[i] {
			t.Errorf("task %d: input mixed up, got %q", id, got.Input)
		}
		if i == 1 {
			if got.State != StateCompleted {
				t.Errorf("task %d: expected %s, got %s", id, StateCompleted, got.State)
			}
			continue
		}
		if got.State != StatePending {
			t.Errorf("task %d: expected %s, got %s", id, StatePending, got.State)
		}
		if got.Result != (Result{}) {
			t.Errorf("task %d: result leaked from another task: %+v", id, got.Result)
		}
	}

	for _, i := range []int{0, 2, 3} {
		if err := m.Process(ids[i]); err != nil {
			t.Fatalf("process: %v", err)
		}
	}

	want := []Result{
		{Words: 1, Runes: 2, Bytes: 2},
		{Words: 3, Runes: 14, Bytes: 14},
		{Words: 2, Runes: 10, Bytes: 19},
		{Words: 5, Runes: 23, Bytes: 23},
	}
	for i, id := range ids {
		got, _ := m.Get(id)
		if got.Result != want[i] {
			t.Errorf("task %d: expected %+v, got %+v", id, want[i], got.Result)
		}
	}

	list := m.List()
	if len(list) != len(inputs) {
		t.Fatalf("expected %d tasks, got %d", len(inputs), len(list))
	}
	for i, tk := range list {
		if tk.ID != ids[i] {
			t.Errorf("list not sorted by ID at position %d: got %d", i, tk.ID)
		}
	}
}

func TestNewManagerDefaultProcessor(t *testing.T) {
	m := NewManager(nil)
	id, _ := m.Create("Hello Go world")

	if err := m.Process(id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := m.Get(id)
	if got.Result.Words != 3 {
		t.Errorf("expected 3 words, got %d", got.Result.Words)
	}
}
