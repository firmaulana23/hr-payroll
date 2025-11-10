package domain

import "time"

// Attendance adalah entitas bisnis inti untuk kehadiran harian
type Attendance struct {
	ID         uint       `json:"id" gorm:"primaryKey" example:"1"`
	EmployeeID uint       `json:"employee_id" example:"1"`
	Date       time.Time  `json:"date" gorm:"uniqueIndex:idx_employee_date" example:"2025-11-10T00:00:00Z"` // Memastikan unik per employee per hari
	Status     string     `json:"status" example:"PRESENT"`                                    // PRESENT, ABSENT, LEAVE
	CheckIn    *time.Time `json:"check_in" example:"2025-11-10T09:00:00Z"`
	CheckOut   *time.Time `json:"check_out" example:"2025-11-10T17:00:00Z"`
	CreatedAt  time.Time  `json:"created_at"`
}

// AttendanceRepository mendefinisikan kontrak operasi data (Port)
type AttendanceRepository interface {
	Save(att *Attendance) error
	Update(att *Attendance) error
	FindByEmployeeAndDate(employeeID uint, date time.Time) (*Attendance, error)
	FindByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]Attendance, error)
}

// AttendanceService mendefinisikan kontrak Use Case
type AttendanceService interface {
	RecordAttendance(att *Attendance) (*Attendance, error)
	RecordCheckout(employeeID uint, checkOutTime time.Time) (*Attendance, error)
	GetAttendanceByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]Attendance, error)
}
