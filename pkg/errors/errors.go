package errors

import (
	"errors"
)

var (
	// Main
	ErrNoPrompt    = errors.New("prompt not provided")
	ErrRecoverable = errors.New("recoverable error")

	// LLM
	ErrUnknownType = errors.New("unknown model type")
)
