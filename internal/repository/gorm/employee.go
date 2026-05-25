package gormrepo

import (
	"context"

	"github.com/danil/org-structure-api/internal/domain"
	"github.com/danil/org-structure-api/internal/repository"
	"gorm.io/gorm"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) repository.EmployeeRepository {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Create(ctx context.Context, e *domain.Employee) error {
	m := employeeModel{
		DepartmentID: e.DepartmentID,
		FullName:     e.FullName,
		Position:     e.Position,
		HiredAt:      e.HiredAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	e.CreatedAt = m.CreatedAt
	return nil
}

func (r *EmployeeRepo) ListByDepartment(ctx context.Context, departmentID uint) ([]domain.Employee, error) {
	var rows []employeeModel
	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Order("full_name asc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.Employee, len(rows))
	for i := range rows {
		out[i] = *toDomainEmployee(&rows[i])
	}
	return out, nil
}

func (r *EmployeeRepo) ReassignFromDepartment(ctx context.Context, fromDeptID, toDeptID uint) error {
	return r.db.WithContext(ctx).
		Model(&employeeModel{}).
		Where("department_id = ?", fromDeptID).
		Update("department_id", toDeptID).Error
}

func (r *EmployeeRepo) DeleteByDepartment(ctx context.Context, departmentID uint) error {
	return r.db.WithContext(ctx).Where("department_id = ?", departmentID).Delete(&employeeModel{}).Error
}

func toDomainEmployee(m *employeeModel) *domain.Employee {
	return &domain.Employee{
		ID:           m.ID,
		DepartmentID: m.DepartmentID,
		FullName:     m.FullName,
		Position:     m.Position,
		HiredAt:      m.HiredAt,
		CreatedAt:    m.CreatedAt,
	}
}
