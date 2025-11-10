package domain

import "time"

// Attendance adalah entitas bisnis inti untuk kehadiran harian
type Attendance struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	EmployeeID uint       `json:"employee_id"`
	Date       time.Time  `json:"date" gorm:"uniqueIndex:idx_employee_date"` // Memastikan unik per employee per hari
	Status     string     `json:"status"`                                    // PRESENT, ABSENT, LEAVE
	CheckIn    *time.Time `json:"check_in"`
	CheckOut   *time.Time `json:"check_out"`
	CreatedAt  time.Time  `json:"created_at"`
}

// AttendanceRepository mendefinisikan kontrak operasi data (Port)
type AttendanceRepository interface {
	Save(att *Attendance) error
	FindByEmployeeAndDate(employeeID uint, date time.Time) (*Attendance, error)
	FindByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]Attendance, error)
}

// AttendanceService mendefinisikan kontrak Use Case
type AttendanceService interface {
	RecordAttendance(att *Attendance) (*Attendance, error)
	GetAttendanceByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]Attendance, error)
}
