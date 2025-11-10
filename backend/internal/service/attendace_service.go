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
	// 1. Validasi: Satu employee hanya boleh satu attendance per hari [cite: 35]
	existingAtt, _ := s.Repo.FindByEmployeeAndDate(att.EmployeeID, att.Date)
	if existingAtt != nil && existingAtt.ID != 0 {
		return nil, errors.New("attendance already recorded for this employee on this date [cite: 35]")
	}

	// 2. Validasi: Status valid
	if att.Status != "PRESENT" && att.Status != "ABSENT" && att.Status != "LEAVE" {
		return nil, errors.New("invalid attendance status: must be PRESENT, ABSENT, or LEAVE [cite: 37]")
	}

	// 3. Validasi: Jika status PRESENT, waktu_datang dan waktu_pulang wajib diisi [cite: 36]
	if att.Status == "PRESENT" {
		if att.CheckIn == nil || att.CheckOut == nil {
			return nil, errors.New("check-in and check-out times are mandatory for PRESENT status [cite: 36]")
		}
	}

	// Simpan ke repository
	if err := s.Repo.Save(att); err != nil {
		return nil, err
	}
	return att, nil
}

// Implementasi GetAttendanceByPeriod
func (s *AttendanceServiceImpl) GetAttendanceByPeriod(employeeID uint, dateFrom time.Time, dateTo time.Time) ([]domain.Attendance, error) {
	return s.Repo.FindByPeriod(employeeID, dateFrom, dateTo)
}
