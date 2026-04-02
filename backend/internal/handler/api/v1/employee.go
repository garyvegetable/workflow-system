package v1

import (
	"log"
	"net/http"
	"strconv"
	"workflow-system/internal/domain/employee"
	"workflow-system/internal/pkg/email"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service   *service.EmployeeService
	emailSvc *email.EmailService
}

func NewEmployeeHandler(service *service.EmployeeService, emailSvc *email.EmailService) *EmployeeHandler {
	return &EmployeeHandler{service: service, emailSvc: emailSvc}
}

func (h *EmployeeHandler) List(c *gin.Context) {
	companyID, _ := strconv.ParseInt(c.Query("company_id"), 10, 64)
	if companyID == 0 {
		companyID = 1
	}

	employees, err := h.service.List(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) Search(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	companyID, _ := strconv.ParseInt(c.Query("company_id"), 10, 64)
	if companyID == 0 {
		companyID = 1
	}

	employees, err := h.service.SearchByName(name, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Level    string `json:"level"`
		CompanyID int64 `json:"company_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plainPassword := req.Password
	emp := &employee.Employee{
		Username: req.Username,
		PasswordHash: req.Password,
		Name: req.Name,
		Email: req.Email,
		Level: req.Level,
		CompanyID: req.CompanyID,
	}

	if err := h.service.Create(emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 发送账号凭证邮件
	go func() {
		if h.emailSvc != nil && req.Email != "" {
			if err := h.emailSvc.SendEmployeeCredentials(req.Email, req.Username, plainPassword); err != nil {
				log.Printf("Failed to send employee credentials email: %v", err)
			}
		}
	}()

	c.JSON(http.StatusCreated, emp)
}

func (h *EmployeeHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	emp, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var emp employee.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emp.ID = id

	if err := h.service.Update(&emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *EmployeeHandler) ListBankAccounts(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	accounts, err := h.service.ListBankAccounts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *EmployeeHandler) CreateBankAccount(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var acc employee.EmployeeBankAccount
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	acc.EmployeeID = id

	if err := h.service.CreateBankAccount(&acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, acc)
}

func (h *EmployeeHandler) UpdateBankAccount(c *gin.Context) {
	aid, _ := strconv.ParseInt(c.Param("aid"), 10, 64)
	var acc employee.EmployeeBankAccount
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	acc.ID = aid

	if err := h.service.UpdateBankAccount(&acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, acc)
}

func (h *EmployeeHandler) DeleteBankAccount(c *gin.Context) {
	aid, _ := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err := h.service.DeleteBankAccount(aid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *EmployeeHandler) SetDepartments(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		DepartmentIDs []int64 `json:"department_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetDepartments(id, req.DepartmentIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "departments updated"})
}

func (h *EmployeeHandler) GetDepartments(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	deptIDs, err := h.service.GetDepartments(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deptIDs)
}
