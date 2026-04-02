package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct{}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h *AuditHandler) List(c *gin.Context) {
	// Mock 数据
	c.JSON(http.StatusOK, []gin.H{
		{"id": 1, "action": "create", "resource_type": "workflow_instance", "resource_id": "1", "created_at": "2024-01-01T10:00:00Z"},
		{"id": 2, "action": "approve", "resource_type": "approval_task", "resource_id": "1", "created_at": "2024-01-01T11:00:00Z"},
	})
}
