package http

import (
	"hr-payroll/internal/delivery/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig menampung semua handler yang dibutuhkan
type RouterConfig struct {
	EmployeeHandler   *handler.EmployeeHandler
	AttendanceHandler *handler.AttendanceHandler
	PayrollHandler    *handler.PayrollHandler
}

// SetupRouter mengkonfigurasi dan mengembalikan router Gin
func SetupRouter(cfg RouterConfig) *gin.Engine {
	router := gin.Default()

	// Endpoint Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Grouping API Version 1
	v1 := router.Group("/api/v1")
	{
		// 1. Employee Management Routes
		v1.POST("/employees", cfg.EmployeeHandler.CreateEmployee)
		v1.GET("/employees", cfg.EmployeeHandler.GetAllEmployees)
		v1.GET("/employees/:id", cfg.EmployeeHandler.GetEmployeeByID)
		// Tambahkan: v1.PUT("/employees/:id", cfg.EmployeeHandler.UpdateEmployee)

		// 2. Attendance Management Routes
		v1.POST("/attendances", cfg.AttendanceHandler.RecordAttendance)
		// Tambahkan: v1.GET("/attendances", cfg.AttendanceHandler.GetAttendanceByPeriod)

		// 3. Payroll Generation Routes
		v1.POST("/payroll/generate", cfg.PayrollHandler.GeneratePayroll)
		v1.GET("/payroll/slips", cfg.PayrollHandler.GetPayrollSlips)
		v1.GET("/payroll/slips/:id", cfg.PayrollHandler.GetPayrollDetail)
	}

	return router
}
