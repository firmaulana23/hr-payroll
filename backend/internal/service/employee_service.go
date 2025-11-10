package service

import "hr-payroll/internal/domain"

// EmployeeServiceImpl mengimplementasikan domain.EmployeeService
type EmployeeServiceImpl struct {
	Repo domain.EmployeeRepository // Dependency pada Interface Repository
}

func NewEmployeeServiceImpl(repo domain.EmployeeRepository) domain.EmployeeService {
	return &EmployeeServiceImpl{Repo: repo}
}

// CreateEmployee implements domain.EmployeeService
func (s *EmployeeServiceImpl) CreateEmployee(emp *domain.Employee) (*domain.Employee, error) {
	// *LOGIKA BISNIS/VALIDASI di sini, jika ada
	// Contoh: memastikan gaji tidak negatif

	if err := s.Repo.Save(emp); err != nil {
		return nil, err
	}
	return emp, nil
}

// GetEmployeeByID implements domain.EmployeeService
func (s *EmployeeServiceImpl) GetEmployeeByID(id uint) (*domain.Employee, error) {
	employee, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err // GORM akan mengembalikan gorm.ErrRecordNotFound jika tidak ditemukan
	}
	return employee, nil
}

// GetAllEmployees implements domain.EmployeeService
func (s *EmployeeServiceImpl) GetAllEmployees() ([]domain.Employee, error) {
	employees, err := s.Repo.FindAll()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// UpdateEmployee implements domain.EmployeeService
func (s *EmployeeServiceImpl) UpdateEmployee(id uint, newEmp *domain.Employee) (*domain.Employee, error) {
	// 1. Cek keberadaan
	existingEmp, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Update field
	existingEmp.Name = newEmp.Name
	existingEmp.BaseSalary = newEmp.BaseSalary
	existingEmp.Allowance = newEmp.Allowance
	existingEmp.Position = newEmp.Position

	// 3. Simpan perubahan
	if err := s.Repo.Update(existingEmp); err != nil {
		return nil, err
	}
	return existingEmp, nil
}
