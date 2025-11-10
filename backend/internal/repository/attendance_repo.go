package repository

import (
	"errors"
	"hr-payroll/internal/domain"
	"time"

	"gorm.io/gorm"
)

// AttendanceGormRepository implements domain.AttendanceRepository
type AttendanceGormRepository struct {
	DB *gorm.DB
}

func NewAttendanceGormRepository(db *gorm.DB) domain.AttendanceRepository {
	return &AttendanceGormRepository{DB: db}
}

// Save implements domain.AttendanceRepository.
func (r *AttendanceGormRepository) Save(att *domain.Attendance) error {
	return r.DB.Create(att).Error
}

// Update implements domain.AttendanceRepository.
func (r *AttendanceGormRepository) Update(att *domain.Attendance) error {
	return r.DB.Model(att).Updates(att).Error
}

// FindByEmployeeAndDate implements domain.AttendanceRepository.
func (r *AttendanceGormRepository) FindByEmployeeAndDate(employeeID uint, date time.Time) (*domain.Attendance, error) {
	var attendance domain.Attendance
	// GORM query to find a record by employee_id and date
	err := r.DB.Where("employee_id = ? AND date = ?", employeeID, date).First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, simplifying service layer error handling
		}
		return nil, err
	}
	return &attendance, nil
}

// FindByPeriod implements domain.AttendanceRepository.
func (r *AttendanceGormRepository) FindByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]domain.Attendance, error) {
	var attendances []domain.Attendance
	// GORM query to find records within a date range for a specific employee
	err := r.DB.Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, dateFrom, dateTo).Find(&attendances).Error
	return attendances, err
}
