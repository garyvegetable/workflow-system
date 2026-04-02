package v1

import (
	"net/http"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemSettingHandler struct {
	service *service.SystemSettingService
}

func NewSystemSettingHandler(service *service.SystemSettingService) *SystemSettingHandler {
	return &SystemSettingHandler{service: service}
}

// List 获取所有设置
func (h *SystemSettingHandler) List(c *gin.Context) {
	settings, err := h.service.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 转为 map 方便前端使用
	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, result)
}

// Get 获取单个设置
func (h *SystemSettingHandler) Get(c *gin.Context) {
	key := c.Param("key")
	value, err := h.service.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

// Set 更新或创建设置
func (h *SystemSettingHandler) Set(c *gin.Context) {
	key := c.Param("key")
	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Set(key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key, "value": req.Value})
}

// Delete 删除设置
func (h *SystemSettingHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	if err := h.service.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
