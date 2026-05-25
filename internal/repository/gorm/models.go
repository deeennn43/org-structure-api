package gormrepo

import "time"

// GORM-модели отражают таблицы БД. Отдельно от domain, чтобы слой БД не «протекал» в HTTP.

type departmentModel struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	ParentID  *uint
	CreatedAt time.Time
}

func (departmentModel) TableName() string { return "departments" }

type employeeModel struct {
	ID           uint `gorm:"primaryKey"`
	DepartmentID uint
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

func (employeeModel) TableName() string { return "employees" }
