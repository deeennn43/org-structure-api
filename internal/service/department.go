package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/danil/org-structure-api/internal/apperrors"
	"github.com/danil/org-structure-api/internal/domain"
	"github.com/danil/org-structure-api/internal/repository"
	"github.com/danil/org-structure-api/internal/validation"
)

// DepartmentService — бизнес-логика подразделений (Single Responsibility).
type DepartmentService struct {
	depts repository.DepartmentRepository
	emps  repository.EmployeeRepository
}

func NewDepartmentService(
	depts repository.DepartmentRepository,
	emps repository.EmployeeRepository,
) *DepartmentService {
	return &DepartmentService{depts: depts, emps: emps}
}

type CreateDepartmentInput struct {
	Name     string
	ParentID *uint
}

func (s *DepartmentService) Create(ctx context.Context, in CreateDepartmentInput) (*domain.Department, error) {
	name, err := validation.RequireNonEmptyString("name", in.Name, 200)
	if err != nil {
		return nil, err
	}
	if in.ParentID != nil {
		ok, err := s.depts.Exists(ctx, *in.ParentID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperrors.ErrNotFound
		}
	}
	dup, err := s.depts.ExistsSiblingName(ctx, in.ParentID, name, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, apperrors.ErrDuplicateName
	}
	d := &domain.Department{Name: name, ParentID: in.ParentID}
	if err := s.depts.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

type UpdateDepartmentInput struct {
	ID       uint
	Name     *string
	ParentID *uint
}

func (s *DepartmentService) Update(ctx context.Context, in UpdateDepartmentInput) (*domain.Department, error) {
	d, err := s.depts.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		name, err := validation.RequireNonEmptyString("name", *in.Name, 200)
		if err != nil {
			return nil, err
		}
		dup, err := s.depts.ExistsSiblingName(ctx, d.ParentID, name, d.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, apperrors.ErrDuplicateName
		}
		d.Name = name
	}
	if in.ParentID != nil {
		if err := s.validateParentChange(ctx, d.ID, in.ParentID); err != nil {
			return nil, err
		}
		dup, err := s.depts.ExistsSiblingName(ctx, in.ParentID, d.Name, d.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, apperrors.ErrDuplicateName
		}
		d.ParentID = in.ParentID
	}
	if err := s.depts.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DepartmentService) validateParentChange(ctx context.Context, deptID uint, newParentID *uint) error {
	if newParentID == nil {
		return nil
	}
	if *newParentID == deptID {
		return fmt.Errorf("%w: cannot be own parent", apperrors.ErrCycle)
	}
	ok, err := s.depts.Exists(ctx, *newParentID)
	if err != nil {
		return err
	}
	if !ok {
		return apperrors.ErrNotFound
	}
	if err := s.ensureNoCycle(ctx, deptID, *newParentID); err != nil {
		return err
	}
	return nil
}

// ensureNoCycle: newParent не должен находиться в поддереве deptID.
func (s *DepartmentService) ensureNoCycle(ctx context.Context, deptID, newParentID uint) error {
	current := newParentID
	for {
		if current == deptID {
			return fmt.Errorf("%w: move into own subtree", apperrors.ErrCycle)
		}
		parent, err := s.depts.GetByID(ctx, current)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				return nil
			}
			return err
		}
		if parent.ParentID == nil {
			return nil
		}
		current = *parent.ParentID
	}
}

type GetDepartmentOptions struct {
	Depth             int
	IncludeEmployees  bool
}

func (s *DepartmentService) GetTree(ctx context.Context, id uint, opt GetDepartmentOptions) (*domain.DepartmentTree, error) {
	d, err := s.depts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.buildTree(ctx, *d, opt.Depth, opt.IncludeEmployees)
}

func (s *DepartmentService) buildTree(ctx context.Context, d domain.Department, depth int, includeEmployees bool) (*domain.DepartmentTree, error) {
	tree := &domain.DepartmentTree{Department: d}
	if includeEmployees {
		emps, err := s.emps.ListByDepartment(ctx, d.ID)
		if err != nil {
			return nil, err
		}
		tree.Employees = emps
	}
	if depth <= 1 {
		return tree, nil
	}
	children, err := s.depts.ListChildren(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	for _, ch := range children {
		sub, err := s.buildTree(ctx, ch, depth-1, includeEmployees)
		if err != nil {
			return nil, err
		}
		tree.Children = append(tree.Children, *sub)
	}
	return tree, nil
}

type DeleteDepartmentInput struct {
	ID                     uint
	Mode                   string
	ReassignToDepartmentID *uint
}

func (s *DepartmentService) Delete(ctx context.Context, in DeleteDepartmentInput) error {
	exists, err := s.depts.Exists(ctx, in.ID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.ErrNotFound
	}
	switch in.Mode {
	case "cascade":
		return s.depts.DeleteCascade(ctx, in.ID)
	case "reassign":
		if in.ReassignToDepartmentID == nil {
			return fmt.Errorf("%w: reassign_to_department_id required", apperrors.ErrValidation)
		}
		ok, err := s.depts.Exists(ctx, *in.ReassignToDepartmentID)
		if err != nil {
			return err
		}
		if !ok {
			return apperrors.ErrNotFound
		}
		children, err := s.depts.ListChildren(ctx, in.ID)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return fmt.Errorf("%w: cannot reassign when department has children", apperrors.ErrConflict)
		}
		if err := s.emps.ReassignFromDepartment(ctx, in.ID, *in.ReassignToDepartmentID); err != nil {
			return err
		}
		return s.depts.Delete(ctx, in.ID)
	default:
		return fmt.Errorf("%w: mode must be cascade or reassign", apperrors.ErrValidation)
	}
}
