package domain

import "time"

// Department — подразделение в организационной структуре.
// ParentID = nil означает корневое подразделение (без родителя).
type Department struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

// DepartmentTree — подразделение с вложенными дочерними и сотрудниками (для GET).
type DepartmentTree struct {
	Department Department        `json:"department"`
	Employees  []Employee        `json:"employees,omitempty"`
	Children   []DepartmentTree  `json:"children,omitempty"`
}
