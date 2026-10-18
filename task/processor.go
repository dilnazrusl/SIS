package task

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Processor interface {
	Process(input string) (Result, error)
}

type TextAnalyzer struct{}

var _ Processor = TextAnalyzer{}

func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("%w: text cannot be empty or whitespace-only", ErrInvalidInput)
	}
	return nil
}

func (TextAnalyzer) Process(input string) (Result, error) {
	if err := ValidateInput(input); err != nil {
		return Result{}, err
	}
	return Result{
		Words: len(strings.Fields(input)),
		Runes: utf8.RuneCountInString(input),
		Bytes: len(input),
	}, nil
}
