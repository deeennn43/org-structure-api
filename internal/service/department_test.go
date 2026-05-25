package service_test

import (
	"context"
	"testing"

	"github.com/deeennn43/org-structure-api/internal/apperrors"
	"github.com/deeennn43/org-structure-api/internal/domain"
	"github.com/deeennn43/org-structure-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeDeptRepo struct {
	byID map[uint]*domain.Department
	next uint
}

func newFakeDeptRepo(seed map[uint]*domain.Department) *fakeDeptRepo {
	by := make(map[uint]*domain.Department, len(seed))
	for k, v := range seed {
		c := *v
		by[k] = &c
	}
	return &fakeDeptRepo{byID: by, next: 100}
}

func (f *fakeDeptRepo) Create(ctx context.Context, d *domain.Department) error {
	f.next++
	d.ID = f.next
	f.byID[d.ID] = d
	return nil
}
func (f *fakeDeptRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	d, ok := f.byID[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	cp := *d
	return &cp, nil
}
func (f *fakeDeptRepo) Update(ctx context.Context, d *domain.Department) error {
	if _, ok := f.byID[d.ID]; !ok {
		return apperrors.ErrNotFound
	}
	cp := *d
	f.byID[d.ID] = &cp
	return nil
}
func (f *fakeDeptRepo) Delete(ctx context.Context, id uint) error { delete(f.byID, id); return nil }
func (f *fakeDeptRepo) Exists(ctx context.Context, id uint) (bool, error) {
	_, ok := f.byID[id]
	return ok, nil
}
func (f *fakeDeptRepo) ExistsSiblingName(ctx context.Context, parentID *uint, name string, excludeID uint) (bool, error) {
	return false, nil
}
func (f *fakeDeptRepo) ListChildren(ctx context.Context, parentID uint) ([]domain.Department, error) {
	return nil, nil
}
func (f *fakeDeptRepo) DeleteCascade(ctx context.Context, id uint) error { return f.Delete(ctx, id) }

type fakeEmpRepo struct{}

func (f *fakeEmpRepo) Create(ctx context.Context, e *domain.Employee) error { return nil }
func (f *fakeEmpRepo) ListByDepartment(ctx context.Context, departmentID uint) ([]domain.Employee, error) {
	return nil, nil
}
func (f *fakeEmpRepo) ReassignFromDepartment(ctx context.Context, fromDeptID, toDeptID uint) error {
	return nil
}
func (f *fakeEmpRepo) DeleteByDepartment(ctx context.Context, departmentID uint) error { return nil }

func TestUpdate_RejectsCycle(t *testing.T) {
	parent := uint(1)
	child := uint(2)

	repo := newFakeDeptRepo(map[uint]*domain.Department{
		1: {ID: 1, Name: "HQ", ParentID: nil},
		2: {ID: 2, Name: "IT", ParentID: &parent},
		3: {ID: 3, Name: "Dev", ParentID: &child},
	})
	svc := service.NewDepartmentService(repo, &fakeEmpRepo{})

	newParent := uint(3)
	_, err := svc.Update(context.Background(), service.UpdateDepartmentInput{
		ID:       2,
		ParentID: &newParent,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrCycle)
}
