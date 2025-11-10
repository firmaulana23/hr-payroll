package handler

import (
	"hr-payroll/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PayrollHandler mengurus endpoint HTTP untuk payroll
type PayrollHandler struct {
	Service domain.PayrollService
}

func NewPayrollHandler(s domain.PayrollService) *PayrollHandler {
	return &PayrollHandler{Service: s}
}

// GeneratePayrollRequest represents the payload to generate payroll
type GeneratePayrollRequest struct {
	EmployeeID uint   `json:"employee_id"`
	Period     string `json:"period"` // expect YYYY-MM-DD (start of month)
}

// GeneratePayroll godoc
// @Summary Generate monthly payroll for an employee
// @Tags Payroll
// @Accept json
// @Produce json
// @Param payload body GeneratePayrollRequest true "Payroll request"
// @Success 201 {object} domain.Payroll
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payroll/generate [post]
func (h *PayrollHandler) GeneratePayroll(c *gin.Context) {
	var req struct {
		EmployeeID uint   `json:"employee_id"`
		Period     string `json:"period"` // expect YYYY-MM-DD (start of month)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// parse period
	period, err := time.Parse("2006-01-02", req.Period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period format. Use YYYY-MM-DD"})
		return
	}

	payroll, err := h.Service.GenerateMonthlyPayroll(req.EmployeeID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payroll)
}

// GetPayrollSlips handles GET /payroll/slips
// GetPayrollSlips godoc
// @Summary List payroll slips
// @Tags Payroll
// @Accept json
// @Produce json
// @Success 200 {array} domain.Payroll
// @Failure 500 {object} map[string]string
// @Router /payroll/slips [get]
func (h *PayrollHandler) GetPayrollSlips(c *gin.Context) {
	slips, err := h.Service.GetPayrollSlips()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payroll slips"})
		return
	}
	c.JSON(http.StatusOK, slips)
}

// GetPayrollDetail handles GET /payroll/slips/:id
// GetPayrollDetail godoc
// @Summary Get payroll detail by ID
// @Tags Payroll
// @Accept json
// @Produce json
// @Param id path int true "Payroll ID"
// @Success 200 {object} domain.Payroll
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payroll/slips/{id} [get]
func (h *PayrollHandler) GetPayrollDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	payroll, err := h.Service.GetPayrollDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payroll detail"})
		return
	}

	c.JSON(http.StatusOK, payroll)
}
