package domain

import "time"

type Department struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type DepartmentTree struct {
	Department Department        `json:"department"`
	Employees  []Employee        `json:"employees,omitempty"`
	Children   []DepartmentTree  `json:"children,omitempty"`
}
