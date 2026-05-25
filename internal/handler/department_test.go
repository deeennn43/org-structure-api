package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danil/org-structure-api/internal/apperrors"
	"github.com/danil/org-structure-api/internal/domain"
	"github.com/danil/org-structure-api/internal/handler"
	"github.com/danil/org-structure-api/internal/service"
	"github.com/stretchr/testify/assert"
)

type stubDeptRepo struct{}

func (s *stubDeptRepo) Create(ctx context.Context, d *domain.Department) error { return nil }
func (s *stubDeptRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	return nil, apperrors.ErrNotFound
}
func (s *stubDeptRepo) Update(ctx context.Context, d *domain.Department) error { return nil }
func (s *stubDeptRepo) Delete(ctx context.Context, id uint) error              { return nil }
func (s *stubDeptRepo) Exists(ctx context.Context, id uint) (bool, error)        { return true, nil }
func (s *stubDeptRepo) ExistsSiblingName(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	return false, nil
}
func (s *stubDeptRepo) ListChildren(ctx context.Context, parentID uint) ([]domain.Department, error) {
	return nil, nil
}
func (s *stubDeptRepo) DeleteCascade(ctx context.Context, id uint) error { return nil }

type stubEmpRepo struct{}

func (s *stubEmpRepo) Create(ctx context.Context, e *domain.Employee) error { return nil }
func (s *stubEmpRepo) ListByDepartment(ctx context.Context, departmentID uint) ([]domain.Employee, error) {
	return nil, nil
}
func (s *stubEmpRepo) ReassignFromDepartment(ctx context.Context, fromDeptID, toDeptID uint) error {
	return nil
}
func (s *stubEmpRepo) DeleteByDepartment(ctx context.Context, departmentID uint) error { return nil }

func TestCreateDepartment_ValidationEmptyName(t *testing.T) {
	deptRepo := &stubDeptRepo{}
	empRepo := &stubEmpRepo{}
	deptSvc := service.NewDepartmentService(deptRepo, empRepo)
	empSvc := service.NewEmployeeService(deptRepo, empRepo)
	h := handler.NewDepartmentHandler(deptSvc, empSvc)
	router := handler.NewRouter(h)

	body := bytes.NewBufferString(`{"name":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/departments/", body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
