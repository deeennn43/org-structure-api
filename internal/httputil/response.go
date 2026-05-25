package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danil/org-structure-api/internal/apperrors"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode json", "err", err)
	}
}

type ErrorBody struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorBody{Error: err.Error()})
	case errors.Is(err, apperrors.ErrConflict), errors.Is(err, apperrors.ErrDuplicateName), errors.Is(err, apperrors.ErrCycle):
		WriteJSON(w, http.StatusConflict, ErrorBody{Error: err.Error()})
	case errors.Is(err, apperrors.ErrValidation):
		WriteJSON(w, http.StatusBadRequest, ErrorBody{Error: err.Error()})
	default:
		slog.Error("internal error", "err", err)
		WriteJSON(w, http.StatusInternalServerError, ErrorBody{Error: "internal server error"})
	}
}

func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
