package handler

import (
	"hr-payroll/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	Service domain.AttendanceService
}

func NewAttendanceHandler(s domain.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{Service: s}
}

// RecordAttendance handles POST /attendances
// @Summary Record daily attendance
// @Tags Attendances
// @Accept json
// @Produce json
// @Param attendance body domain.Attendance true "Attendance object"
// @Success 201 {object} domain.Attendance
// @Router /attendances [post]
func (h *AttendanceHandler) RecordAttendance(c *gin.Context) {
	var req domain.Attendance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Pastikan Date di-set ke awal hari untuk validasi unik yang benar
	req.Date = time.Date(req.Date.Year(), req.Date.Month(), req.Date.Day(), 0, 0, 0, 0, time.UTC)

	attendance, err := h.Service.RecordAttendance(&req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attendance)
}

// RecordCheckout handles PUT /attendances/checkout
// @Summary Record checkout for an employee
// @Tags Attendances
// @Accept json
// @Produce json
// @Param checkout body domain.Attendance true "Checkout object"
// @Success 200 {object} domain.Attendance
// @Router /attendances/checkout [put]
func (h *AttendanceHandler) RecordCheckout(c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	now := time.Now()
	attendance, err := h.Service.RecordCheckout(req.EmployeeID, now)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendance)
}

// GetAttendanceByPeriod handles GET /attendances
// @Summary Get attendance records by period for an employee
// @Tags Attendances
// @Accept json
// @Produce json
// @Param employee_id query int true "Employee ID"
// @Param from query string true "From date (YYYY-MM-DD)"
// @Param to query string true "To date (YYYY-MM-DD)"
// @Success 200 {array} domain.Attendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendances [get]
func (h *AttendanceHandler) GetAttendanceByPeriod(c *gin.Context) {
	employeeIDStr := c.Query("employee_id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee_id format"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' date format, use YYYY-MM-DD"})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' date format, use YYYY-MM-DD"})
		return
	}

	attendances, err := h.Service.GetAttendanceByPeriod(uint(employeeID), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendances"})
		return
	}

	c.JSON(http.StatusOK, attendances)
}
