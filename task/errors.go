package task

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid state transition")
)
