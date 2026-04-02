package v1

import (
	"net/http"
	"strconv"
	"workflow-system/internal/domain/expense_category"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type ExpenseCategoryHandler struct {
	service *service.ExpenseCategoryService
}

func NewExpenseCategoryHandler(service *service.ExpenseCategoryService) *ExpenseCategoryHandler {
	return &ExpenseCategoryHandler{service: service}
}

func (h *ExpenseCategoryHandler) List(c *gin.Context) {
	companyID, _ := strconv.ParseInt(c.Query("company_id"), 10, 64)
	if companyID == 0 {
		companyID = 1
	}

	categories, err := h.service.GetTree(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *ExpenseCategoryHandler) Create(c *gin.Context) {
	var cat expense_category.ExpenseCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *ExpenseCategoryHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cat, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *ExpenseCategoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var cat expense_category.ExpenseCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id

	if err := h.service.Update(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *ExpenseCategoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
