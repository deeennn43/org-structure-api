package validation_test

import (
	"testing"

	"github.com/deeennn43/org-structure-api/internal/apperrors"
	"github.com/deeennn43/org-structure-api/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireNonEmptyString_TrimsSpaces(t *testing.T) {
	v, err := validation.RequireNonEmptyString("name", "  Backend  ", 200)
	require.NoError(t, err)
	assert.Equal(t, "Backend", v)
}

func TestRequireNonEmptyString_RejectsEmpty(t *testing.T) {
	_, err := validation.RequireNonEmptyString("name", "   ", 200)
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
}
