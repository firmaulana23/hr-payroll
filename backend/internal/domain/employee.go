package domain

import "time"

// Employee adalah entitas bisnis inti
type Employee struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	BaseSalary float64   `json:"base_salary"`
	Allowance  float64   `json:"allowance"`
	Position   string    `json:"position"`
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
