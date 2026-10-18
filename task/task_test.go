package task

import (
	"errors"
	"testing"
)

func TestTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    State
		to      State
		wantErr bool
	}{
		{"pending to processing", StatePending, StateProcessing, false},
		{"processing to completed", StateProcessing, StateCompleted, false},
		{"processing to failed", StateProcessing, StateFailed, false},
		{"pending to completed", StatePending, StateCompleted, true},
		{"pending to failed", StatePending, StateFailed, true},
		{"processing to pending", StateProcessing, StatePending, true},
		{"processing to processing", StateProcessing, StateProcessing, true},
		{"completed to processing", StateCompleted, StateProcessing, true},
		{"completed to pending", StateCompleted, StatePending, true},
		{"failed to processing", StateFailed, StateProcessing, true},
		{"unknown state", State("Unknown"), StateProcessing, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk := &Task{State: tt.from}
			err := tk.transition(tt.to)

			if tt.wantErr {
				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("expected ErrInvalidTransition, got %v", err)
				}
				if tk.State != tt.from {
					t.Errorf("state changed on rejected transition: %s", tk.State)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tk.State != tt.to {
				t.Errorf("expected state %s, got %s", tt.to, tk.State)
			}
		})
	}
}

func TestTextAnalyzerInvalidInput(t *testing.T) {
	_, err := TextAnalyzer{}.Process("   ")
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}
