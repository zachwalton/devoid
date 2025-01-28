package errors

import (
	"errors"
	"fmt"
)

var (
	// Main
	ErrNoPrompt    = errors.New("prompt not provided")
	ErrRecoverable = errors.New("recoverable error")

	// LLM
	ErrUnknownType = errors.New("unknown model type")

	// Initial stage
	ErrNoProjectPath     = errors.New("project path was not set")
	ErrAbsolutePath      = errors.New("failed to resolve absolute path")
	ErrProjectPathEscape = fmt.Errorf("%w: %s", ErrRecoverable)
)
