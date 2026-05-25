package domain

import "time"

// Employee — сотрудник, привязанный к одному подразделению.
type Employee struct {
	ID           uint       `json:"id"`
	DepartmentID uint       `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
