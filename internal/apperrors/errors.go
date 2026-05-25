package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrValidation     = errors.New("validation")
	ErrCycle          = errors.New("department cycle")
	ErrDuplicateName  = errors.New("duplicate department name")
)
