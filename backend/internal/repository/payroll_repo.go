package repository

import (
	"errors"
	"hr-payroll/internal/domain"
	"time"

	"gorm.io/gorm"
)

// PayrollGormRepository implements domain.PayrollRepository
type PayrollGormRepository struct {
	DB *gorm.DB
}

func NewPayrollGormRepository(db *gorm.DB) domain.PayrollRepository {
	return &PayrollGormRepository{DB: db}
}

// Save implements domain.PayrollRepository.
func (r *PayrollGormRepository) Save(payroll *domain.Payroll) error {
	return r.DB.Create(payroll).Error
}

// FindByEmployeeAndPeriod implements domain.PayrollRepository.
func (r *PayrollGormRepository) FindByEmployeeAndPeriod(employeeID uint, period time.Time) (*domain.Payroll, error) {
	var payroll domain.Payroll
	// GORM query to check for existing payroll based on unique constraint
	err := r.DB.Where("employee_id = ? AND period = ?", employeeID, period).First(&payroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payroll, nil
}

// FindAll implements domain.PayrollRepository.
func (r *PayrollGormRepository) FindAll() ([]domain.Payroll, error) {
	var payrolls []domain.Payroll
	err := r.DB.Find(&payrolls).Error
	return payrolls, err
}

// FindByID implements domain.PayrollRepository.
func (r *PayrollGormRepository) FindByID(id uint) (*domain.Payroll, error) {
	var payroll domain.Payroll
	err := r.DB.First(&payroll, id).Error
	if err != nil {
		return nil, err
	}
	return &payroll, nil
}
