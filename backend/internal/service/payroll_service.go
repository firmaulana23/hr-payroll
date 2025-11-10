package service

import (
	"errors"
	"hr-payroll/internal/domain"
	"time"
)

type PayrollServiceImpl struct {
	EmpRepo domain.EmployeeRepository
	AttRepo domain.AttendanceRepository
	PayRepo domain.PayrollRepository
}

func NewPayrollServiceImpl(er domain.EmployeeRepository, ar domain.AttendanceRepository, pr domain.PayrollRepository) domain.PayrollService {
	return &PayrollServiceImpl{EmpRepo: er, AttRepo: ar, PayRepo: pr}
}

// GenerateMonthlyPayroll implements domain.PayrollService
func (s *PayrollServiceImpl) GenerateMonthlyPayroll(employeeID uint, period time.Time) (*domain.Payroll, error) {
	// 1. Validasi Unik: Payroll untuk kombinasi employee_id + period hanya boleh satu [cite: 42]
	existingPayroll, _ := s.PayRepo.FindByEmployeeAndPeriod(employeeID, period)
	if existingPayroll != nil && existingPayroll.ID != 0 {
		return nil, errors.New("payroll already generated for this employee and period [cite: 42]")
	}

	// 2. Ambil data Employee
	employee, err := s.EmpRepo.FindByID(employeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Tentukan periode attendance (Asumsi: sebulan penuh sebelum 'period')
	dateFrom := time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0).Add(-time.Second) // Akhir bulan

	// 3. Ambil data Attendance dan hitung total absent
	attendances, _ := s.AttRepo.FindByPeriod(employeeID, dateFrom, dateTo)
	totalAbsent := 0
	for _, att := range attendances {
		if att.Status == "ABSENT" {
			totalAbsent++
		}
	}

	// 4. Hitung Deduction dan Take Home Pay [cite: 41]
	const workingDaysInMonth = 22
	dailySalary := employee.BaseSalary / workingDaysInMonth
	absenceDeduction := dailySalary * float64(totalAbsent)                     // (base_salary / 22) * total_absent [cite: 41]
	takeHomePay := employee.BaseSalary + employee.Allowance - absenceDeduction // base_salary + allowance - absence_deduction [cite: 41]

	// 5. Buat dan simpan entitas Payroll
	payroll := &domain.Payroll{
		EmployeeID:       employeeID,
		Period:           period,
		BaseSalary:       employee.BaseSalary,
		Allowance:        employee.Allowance,
		TotalAbsent:      totalAbsent,
		AbsenceDeduction: absenceDeduction,
		TakeHomePay:      takeHomePay,
		GeneratedAt:      time.Now(),
	}

	if err := s.PayRepo.Save(payroll); err != nil {
		return nil, err
	}
	return payroll, nil
}

// GetPayrollSlips implements domain.PayrollService
func (s *PayrollServiceImpl) GetPayrollSlips() ([]domain.Payroll, error) {
	// Memanggil repository untuk mengambil semua slip gaji
	payrolls, err := s.PayRepo.FindAll()
	if err != nil {
		// Logika penanganan error khusus bisa ditambahkan di sini
		return nil, err
	}
	return payrolls, nil
}

// GetPayrollDetail implements domain.PayrollService
func (s *PayrollServiceImpl) GetPayrollDetail(id uint) (*domain.Payroll, error) {
	// Memanggil repository untuk mengambil detail slip gaji berdasarkan ID
	payroll, err := s.PayRepo.FindByID(id)
	if err != nil {
		// Asumsi GORM/Repo mengembalikan error spesifik jika tidak ditemukan
		return nil, errors.New("payroll slip not found or database error")
	}
	return payroll, nil
}
