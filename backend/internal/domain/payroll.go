package domain

import "time"

// Payroll adalah entitas bisnis inti untuk slip gaji bulanan
type Payroll struct {
	ID               uint      `json:"id" gorm:"primaryKey" example:"1"`
	EmployeeID       uint      `json:"employee_id" gorm:"uniqueIndex:idx_employee_period" example:"1"`
	Period           time.Time `json:"period" gorm:"uniqueIndex:idx_employee_period" example:"2025-11-01T00:00:00Z"` // Biasanya awal bulan
	BaseSalary       float64   `json:"base_salary" example:"50000"`
	Allowance        float64   `json:"allowance" example:"5000"`
	TotalAbsent      int       `json:"total_absent" example:"2"`
	AbsenceDeduction float64   `json:"absence_deduction" example:"1000"`
	TakeHomePay      float64   `json:"take_home_pay" example:"54000"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// PayrollRepository mendefinisikan kontrak operasi data (Port)
type PayrollRepository interface {
	Save(payroll *Payroll) error
	FindByEmployeeAndPeriod(employeeID uint, period time.Time) (*Payroll, error)
	FindAll() ([]Payroll, error)
	FindByID(id uint) (*Payroll, error)
}

// PayrollService mendefinisikan kontrak Use Case
type PayrollService interface {
	GenerateMonthlyPayroll(employeeID uint, period time.Time) (*Payroll, error)
	GetPayrollSlips() ([]Payroll, error)
	GetPayrollDetail(id uint) (*Payroll, error)
}
