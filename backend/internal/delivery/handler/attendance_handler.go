package handler

import (
	"hr-payroll/internal/domain"
	"net/http"
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

// Tambahkan handler untuk GetAttendanceByPeriod (dengan filtering)
// ...
