package handler

import (
	"hr-payroll/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EmployeeHandler mengurus endpoint HTTP
type EmployeeHandler struct {
	Service domain.EmployeeService // Dependency pada Interface Service
}

func NewEmployeeHandler(s domain.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{Service: s}
}

// CreateEmployee handles POST /employees
// CreateEmployee godoc
// @Summary Create a new employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body domain.Employee true "Employee object"
// @Success 201 {object} domain.Employee
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req domain.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Panggil Service
	employee, err := h.Service.CreateEmployee(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// GetEmployeeByID handles GET /employees/:id
// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} domain.Employee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Panggil Service
	employee, err := h.Service.GetEmployeeByID(uint(id))
	if err != nil {
		// Logika sederhana untuk Not Found
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// GetAllEmployees handles GET /employees
// GetAllEmployees godoc
// @Summary List all employees
// @Tags Employees
// @Accept json
// @Produce json
// @Success 200 {array} domain.Employee
// @Failure 500 {object} map[string]string
// @Router /employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	employees, err := h.Service.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// UpdateEmployee handles PUT /employees/:id
// @Summary Update an existing employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body domain.Employee true "Employee object"
// @Success 200 {object} domain.Employee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req domain.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updatedEmployee, err := h.Service.UpdateEmployee(uint(id), &req)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, updatedEmployee)
}
