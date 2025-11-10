package service

import (
	"errors"
	"hr-payroll/internal/domain"
	"time"
)

type AttendanceServiceImpl struct {
	Repo domain.AttendanceRepository
}

func NewAttendanceServiceImpl(repo domain.AttendanceRepository) domain.AttendanceService {
	return &AttendanceServiceImpl{Repo: repo}
}

// RecordAttendance implements domain.AttendanceService
func (s *AttendanceServiceImpl) RecordAttendance(att *domain.Attendance) (*domain.Attendance, error) {
	// 1. Cek apakah sudah ada absensi untuk employee dan tanggal ini
	existingAtt, _ := s.Repo.FindByEmployeeAndDate(att.EmployeeID, att.Date)

	if existingAtt != nil && existingAtt.ID != 0 {
		return nil, errors.New("attendance already recorded for this employee on this date")
	}

	// Record does not exist. This is a new attendance record (check-in, absent, or leave).
	// 2. Validasi: Status valid
	if att.Status != "PRESENT" && att.Status != "ABSENT" && att.Status != "LEAVE" {
		return nil, errors.New("invalid attendance status: must be PRESENT, ABSENT, or LEAVE")
	}

	// 3. Validasi: Jika status PRESENT, waktu_datang wajib diisi
	if att.Status == "PRESENT" {
		if att.CheckIn == nil {
			return nil, errors.New("check-in time is mandatory for PRESENT status")
		}
	}

	// Simpan ke repository
	if err := s.Repo.Save(att); err != nil {
		return nil, err
	}
	return att, nil
}

// RecordCheckout implements domain.AttendanceService
func (s *AttendanceServiceImpl) RecordCheckout(employeeID uint, checkOutTime time.Time) (*domain.Attendance, error) {
	// 1. Find today's attendance record for the employee
	today := time.Now()
	normalizedDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	existingAtt, err := s.Repo.FindByEmployeeAndDate(employeeID, normalizedDate)
	if err != nil {
		return nil, err
	}
	if existingAtt == nil {
		return nil, errors.New("no check-in record found for today")
	}

	// 2. Check if already checked out
	if existingAtt.CheckOut != nil {
		return nil, errors.New("already checked out for today")
	}

	// 3. Update checkout time
	existingAtt.CheckOut = &checkOutTime

	// 4. Save updated record
	if err := s.Repo.Update(existingAtt); err != nil {
		return nil, err
	}

	return existingAtt, nil
}

// Implementasi GetAttendanceByPeriod
func (s *AttendanceServiceImpl) GetAttendanceByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]domain.Attendance, error) {
	return s.Repo.FindByPeriod(employeeID, dateFrom, dateTo)
}
