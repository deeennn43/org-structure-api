package repository

import (
	"context"

	"github.com/deeennn43/org-structure-api/internal/domain"
)

type EmployeeRepository interface {
	Create(ctx context.Context, e *domain.Employee) error
	ListByDepartment(ctx context.Context, departmentID uint) ([]domain.Employee, error)
	ReassignFromDepartment(ctx context.Context, fromDeptID, toDeptID uint) error
	DeleteByDepartment(ctx context.Context, departmentID uint) error
}
