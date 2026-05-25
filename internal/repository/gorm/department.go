package gormrepo

import (
	"context"
	"errors"

	"github.com/deeennn43/org-structure-api/internal/apperrors"
	"github.com/deeennn43/org-structure-api/internal/domain"
	"github.com/deeennn43/org-structure-api/internal/repository"
	"gorm.io/gorm"
)

type DepartmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) repository.DepartmentRepository {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) Create(ctx context.Context, d *domain.Department) error {
	m := departmentModel{Name: d.Name, ParentID: d.ParentID}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	return nil
}

func (r *DepartmentRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	var m departmentModel
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomainDepartment(&m), nil
}

func (r *DepartmentRepo) Update(ctx context.Context, d *domain.Department) error {
	res := r.db.WithContext(ctx).Model(&departmentModel{}).Where("id = ?", d.ID).Updates(map[string]any{
		"name":      d.Name,
		"parent_id": d.ParentID,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *DepartmentRepo) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&departmentModel{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *DepartmentRepo) DeleteCascade(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		sub := &DepartmentRepo{db: tx}
		children, err := sub.ListChildren(ctx, id)
		if err != nil {
			return err
		}
		for _, ch := range children {
			if err := sub.DeleteCascade(ctx, ch.ID); err != nil {
				return err
			}
		}
		if err := NewEmployeeRepo(tx).DeleteByDepartment(ctx, id); err != nil {
			return err
		}
		return sub.Delete(ctx, id)
	})
}

func (r *DepartmentRepo) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&departmentModel{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *DepartmentRepo) ExistsSiblingName(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).Model(&departmentModel{}).Where("name = ?", name)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DepartmentRepo) ListChildren(ctx context.Context, parentID uint) ([]domain.Department, error) {
	var rows []departmentModel
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("name asc").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.Department, len(rows))
	for i := range rows {
		out[i] = *toDomainDepartment(&rows[i])
	}
	return out, nil
}

func toDomainDepartment(m *departmentModel) *domain.Department {
	return &domain.Department{
		ID:        m.ID,
		Name:      m.Name,
		ParentID:  m.ParentID,
		CreatedAt: m.CreatedAt,
	}
}
