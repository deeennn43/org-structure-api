package apperrors

import "errors"

// HTTP-слой переводит в нужные коды ответа (404, 409, 400).
var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrValidation     = errors.New("validation")
	ErrCycle          = errors.New("department cycle")
	ErrDuplicateName  = errors.New("duplicate department name")
)
