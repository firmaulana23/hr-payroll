package domain

import "time"

// Payroll adalah entitas bisnis inti untuk slip gaji bulanan
type Payroll struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	EmployeeID       uint      `json:"employee_id" gorm:"uniqueIndex:idx_employee_period"`
	Period           time.Time `json:"period" gorm:"uniqueIndex:idx_employee_period"` // Biasanya awal bulan
	BaseSalary       float64   `json:"base_salary"`
	Allowance        float64   `json:"allowance"`
	TotalAbsent      int       `json:"total_absent"`
	AbsenceDeduction float64   `json:"absence_deduction"`
	TakeHomePay      float64   `json:"take_home_pay"`
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
