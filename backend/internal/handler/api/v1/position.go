package v1

import (
	"net/http"
	"strconv"
	"workflow-system/internal/domain/position"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	service *service.PositionService
}

func NewPositionHandler(service *service.PositionService) *PositionHandler {
	return &PositionHandler{service: service}
}

func (h *PositionHandler) List(c *gin.Context) {
	companyID, _ := strconv.ParseInt(c.Query("company_id"), 10, 64)
	if companyID == 0 {
		companyID = 1
	}

	positions, err := h.service.List(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, positions)
}

func (h *PositionHandler) Create(c *gin.Context) {
	var pos position.Position
	if err := c.ShouldBindJSON(&pos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 JWT token 中的 company_id
	companyID := c.GetInt64("company_id")
	if companyID == 0 {
		companyID = 1
	}
	pos.CompanyID = companyID

	if err := h.service.Create(&pos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pos)
}

func (h *PositionHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pos, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, pos)
}

func (h *PositionHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var pos position.Position
	if err := c.ShouldBindJSON(&pos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pos.ID = id

	if err := h.service.Update(&pos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pos)
}

func (h *PositionHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
