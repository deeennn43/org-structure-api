package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/danil/org-structure-api/internal/apperrors"
	"github.com/danil/org-structure-api/internal/config"
	"github.com/danil/org-structure-api/internal/httputil"
	"github.com/danil/org-structure-api/internal/service"
)

type DepartmentHandler struct {
	depts *service.DepartmentService
	emps  *service.EmployeeService
}

func NewDepartmentHandler(depts *service.DepartmentService, emps *service.EmployeeService) *DepartmentHandler {
	return &DepartmentHandler{depts: depts, emps: emps}
}

type createDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createDepartmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, errors.Join(apperrors.ErrValidation, err))
		return
	}
	d, err := h.depts.Create(r.Context(), service.CreateDepartmentInput{
		Name:     req.Name,
		ParentID: req.ParentID,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, d)
}

func (h *DepartmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := config.ParseUintID(r.PathValue("id"))
	if err != nil {
		httputil.WriteError(w, apperrors.ErrValidation)
		return
	}
	depth := parseIntDefault(r.URL.Query().Get("depth"), 1)
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}
	includeEmployees := parseBoolDefault(r.URL.Query().Get("include_employees"), true)

	tree, err := h.depts.GetTree(r.Context(), id, service.GetDepartmentOptions{
		Depth:            depth,
		IncludeEmployees: includeEmployees,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, tree)
}

type patchDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *uint   `json:"parent_id"`
}

func (h *DepartmentHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := config.ParseUintID(r.PathValue("id"))
	if err != nil {
		httputil.WriteError(w, apperrors.ErrValidation)
		return
	}
	var req patchDepartmentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, errors.Join(apperrors.ErrValidation, err))
		return
	}
	d, err := h.depts.Update(r.Context(), service.UpdateDepartmentInput{
		ID:       id,
		Name:     req.Name,
		ParentID: req.ParentID,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, d)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := config.ParseUintID(r.PathValue("id"))
	if err != nil {
		httputil.WriteError(w, apperrors.ErrValidation)
		return
	}
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "cascade"
	}
	var reassignID *uint
	if v := r.URL.Query().Get("reassign_to_department_id"); v != "" {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			httputil.WriteError(w, apperrors.ErrValidation)
			return
		}
		u := uint(n)
		reassignID = &u
	}
	if err := h.depts.Delete(r.Context(), service.DeleteDepartmentInput{
		ID:                     id,
		Mode:                   mode,
		ReassignToDepartmentID: reassignID,
	}); err != nil {
		httputil.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

func (h *DepartmentHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	deptID, err := config.ParseUintID(r.PathValue("id"))
	if err != nil {
		httputil.WriteError(w, apperrors.ErrValidation)
		return
	}
	var req createEmployeeRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, errors.Join(apperrors.ErrValidation, err))
		return
	}
	e, err := h.emps.Create(r.Context(), service.CreateEmployeeInput{
		DepartmentID: deptID,
		FullName:     req.FullName,
		Position:     req.Position,
		HiredAt:      req.HiredAt,
	})
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, e)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func parseBoolDefault(s string, def bool) bool {
	if s == "" {
		return def
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}
