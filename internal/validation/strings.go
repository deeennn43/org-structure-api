package validation

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/deeennn43/org-structure-api/internal/apperrors"
)

func NormalizeName(s string) string {
	return strings.TrimSpace(s)
}

func RequireNonEmptyString(field, value string, max int) (string, error) {
	v := NormalizeName(value)
	if v == "" {
		return "", fmt.Errorf("%w: %s must not be empty", apperrors.ErrValidation, field)
	}
	if utf8.RuneCountInString(v) > max {
		return "", fmt.Errorf("%w: %s max length is %d", apperrors.ErrValidation, field, max)
	}
	return v, nil
}
