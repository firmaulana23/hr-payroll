package database

import (
	"fmt"
	"hr-payroll/config"
	"hr-payroll/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate schema (ONLY for development!)
	// db.AutoMigrate(&domain.Employee{}, &domain.Attendance{}, &domain.Payroll{})
	db.AutoMigrate(&domain.Employee{})

	return db
}
