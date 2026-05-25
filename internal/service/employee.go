package service

import (
	"context"
	"time"

	"github.com/deeennn43/org-structure-api/internal/apperrors"
	"github.com/deeennn43/org-structure-api/internal/domain"
	"github.com/deeennn43/org-structure-api/internal/repository"
	"github.com/deeennn43/org-structure-api/internal/validation"
)

type EmployeeService struct {
	depts repository.DepartmentRepository
	emps  repository.EmployeeRepository
}

func NewEmployeeService(
	depts repository.DepartmentRepository,
	emps repository.EmployeeRepository,
) *EmployeeService {
	return &EmployeeService{depts: depts, emps: emps}
}

type CreateEmployeeInput struct {
	DepartmentID uint
	FullName     string
	Position     string
	HiredAt      *time.Time
}

func (s *EmployeeService) Create(ctx context.Context, in CreateEmployeeInput) (*domain.Employee, error) {
	ok, err := s.depts.Exists(ctx, in.DepartmentID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	fullName, err := validation.RequireNonEmptyString("full_name", in.FullName, 200)
	if err != nil {
		return nil, err
	}
	position, err := validation.RequireNonEmptyString("position", in.Position, 200)
	if err != nil {
		return nil, err
	}
	e := &domain.Employee{
		DepartmentID: in.DepartmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      in.HiredAt,
	}
	if err := s.emps.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}
