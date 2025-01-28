package errors

import "errors"

var (
	// Main
	ErrNoPrompt = errors.New("prompt not provided")

	// LLM
	ErrUnknownType = errors.New("unknown model type")
)
