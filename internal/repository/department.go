package repository

import (
	"context"

	"github.com/deeennn43/org-structure-api/internal/domain"
)

type DepartmentRepository interface {
	Create(ctx context.Context, d *domain.Department) error
	GetByID(ctx context.Context, id uint) (*domain.Department, error)
	Update(ctx context.Context, d *domain.Department) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
	ExistsSiblingName(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error)
	ListChildren(ctx context.Context, parentID uint) ([]domain.Department, error)
	DeleteCascade(ctx context.Context, id uint) error
}
