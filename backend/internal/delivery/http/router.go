package http

import (
	"hr-payroll/internal/delivery/handler"
	"time"

	"github.com/gin-contrib/cors"
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
func SetupRouter(router *gin.Engine, cfg RouterConfig) {
	// Configure CORS using gin-contrib/cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Grouping API Version 1
	v1 := router.Group("/api/v1")
	{
		// 1. Employee Management Routes
		v1.POST("/employees", cfg.EmployeeHandler.CreateEmployee)
		v1.GET("/employees", cfg.EmployeeHandler.GetAllEmployees)
		v1.GET("/employees/:id", cfg.EmployeeHandler.GetEmployeeByID)
		v1.PUT("/employees/:id", cfg.EmployeeHandler.UpdateEmployee)

		// 2. Attendance Management Routes
		v1.POST("/attendances", cfg.AttendanceHandler.RecordAttendance)
		v1.PUT("/attendances/checkout", cfg.AttendanceHandler.RecordCheckout)
		v1.GET("/attendances", cfg.AttendanceHandler.GetAttendanceByPeriod)

		// 3. Payroll Generation Routes
		v1.POST("/payroll/generate", cfg.PayrollHandler.GeneratePayroll)
		v1.GET("/payroll/slips", cfg.PayrollHandler.GetPayrollSlips)
		v1.GET("/payroll/slips/:id", cfg.PayrollHandler.GetPayrollDetail)
	}

}
