package domain

import "time"

// Employee adalah entitas bisnis inti
type Employee struct {
	ID         uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name       string    `json:"name" example:"John Doe"`
	BaseSalary float64   `json:"base_salary" example:"50000"`
	Allowance  float64   `json:"allowance" example:"5000"`
	Position   string    `json:"position" example:"Software Engineer"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EmployeeRepository mendefinisikan kontrak operasi data (Port)
type EmployeeRepository interface {
	Save(emp *Employee) error
	FindByID(id uint) (*Employee, error)
	FindAll() ([]Employee, error)
	Update(emp *Employee) error
}

// EmployeeService mendefinisikan kontrak Use Case
type EmployeeService interface {
	CreateEmployee(emp *Employee) (*Employee, error)
	GetEmployeeByID(id uint) (*Employee, error)
	GetAllEmployees() ([]Employee, error)
	UpdateEmployee(id uint, emp *Employee) (*Employee, error)
}
