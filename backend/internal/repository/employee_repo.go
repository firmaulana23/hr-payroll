package repository

import (
	"hr-payroll/internal/domain"

	"gorm.io/gorm"
)

// EmployeeGormRepository adalah adapter yang mengimplementasikan domain.EmployeeRepository
type EmployeeGormRepository struct {
	DB *gorm.DB
}

func NewEmployeeGormRepository(db *gorm.DB) domain.EmployeeRepository {
	return &EmployeeGormRepository{DB: db}
}

// Save implements domain.EmployeeRepository.
func (r *EmployeeGormRepository) Save(emp *domain.Employee) error {
	// Menciptakan karyawan baru
	return r.DB.Create(emp).Error
}

// FindByID implements domain.EmployeeRepository.
func (r *EmployeeGormRepository) FindByID(id uint) (*domain.Employee, error) {
	var employee domain.Employee
	// Mencari karyawan berdasarkan ID
	if err := r.DB.First(&employee, id).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindAll implements domain.EmployeeRepository.
func (r *EmployeeGormRepository) FindAll() ([]domain.Employee, error) {
	var employees []domain.Employee
	// Mengambil semua karyawan
	if err := r.DB.Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

// Update implements domain.EmployeeRepository.
func (r *EmployeeGormRepository) Update(emp *domain.Employee) error {
	// Memperbarui data karyawan
	return r.DB.Save(emp).Error
}
