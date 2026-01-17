package engine

import "errors"

var (
	ErrInvalidLanguage = errors.New("invalid language")
	ErrSourceTooLarge  = errors.New("source too large")
	ErrTimeout         = errors.New("time limit exceeded")
)
