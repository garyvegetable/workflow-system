package v1

import (
	"net/http"
	"strconv"
	"workflow-system/internal/repository"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	service     *service.ApprovalService
	employeeRepo *repository.EmployeeRepository
}

func NewApprovalHandler(service *service.ApprovalService, employeeRepo *repository.EmployeeRepository) *ApprovalHandler {
	return &ApprovalHandler{
		service:     service,
		employeeRepo: employeeRepo,
	}
}

func (h *ApprovalHandler) ListPending(c *gin.Context) {
	// 从 JWT 获取当前用户ID和公司ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// 直接从 JWT claims 获取 company_id（middleware 已设置）
	companyID := c.GetInt64("company_id")
	if companyID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "公司信息无效"})
		return
	}

	tasks, err := h.service.ListPending(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *ApprovalHandler) ListHandled(c *gin.Context) {
	// 从 JWT 获取当前用户ID和公司ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// 直接从 JWT claims 获取 company_id（middleware 已设置）
	companyID := c.GetInt64("company_id")
	if companyID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "公司信息无效"})
		return
	}

	tasks, err := h.service.ListHandled(userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Approve(id, approverID, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Reject(id, approverID, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (h *ApprovalHandler) Transfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		NewAssigneeID int64 `json:"new_assignee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Transfer(id, req.NewAssigneeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transferred"})
}

func (h *ApprovalHandler) GetHistory(c *gin.Context) {
	instanceID, _ := strconv.ParseInt(c.Query("instance_id"), 10, 64)
	history, err := h.service.GetHistory(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// BatchApprove 批量审批
func (h *ApprovalHandler) BatchApprove(c *gin.Context) {
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		TaskIDs []int64 `json:"task_ids"`
		Comment string  `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.TaskIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要审批的任务"})
		return
	}

	successCount := 0
	failedCount := 0
	for _, taskID := range req.TaskIDs {
		if err := h.service.Approve(taskID, approverID, req.Comment); err != nil {
			failedCount++
		} else {
			successCount++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "批量审批完成",
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}

// BatchReject 批量驳回
func (h *ApprovalHandler) BatchReject(c *gin.Context) {
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		TaskIDs []int64 `json:"task_ids"`
		Comment string  `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.TaskIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要驳回的任务"})
		return
	}

	successCount := 0
	failedCount := 0
	for _, taskID := range req.TaskIDs {
		if err := h.service.Reject(taskID, approverID, req.Comment); err != nil {
			failedCount++
		} else {
			successCount++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "批量驳回完成",
		"success_count": successCount,
		"failed_count":  failedCount,
	})
}

// AddApprover 加签：为任务添加审批人
func (h *ApprovalHandler) AddApprover(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		NewApproverID int64 `json:"new_approver_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddApprover(id, req.NewApproverID, approverID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "加签成功"})
}

// RemoveApprover 减签：移除待审批的审批人
func (h *ApprovalHandler) RemoveApprover(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	approverID := c.GetInt64("user_id")
	if approverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		TargetAssigneeID int64 `json:"target_assignee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveApprover(id, req.TargetAssigneeID, approverID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "减签成功"})
}
