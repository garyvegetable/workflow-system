package v1

import (
	"net/http"
	"strconv"
	"workflow-system/internal/domain/department"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	companyID, _ := strconv.ParseInt(c.Query("company_id"), 10, 64)
	if companyID == 0 {
		companyID = 1 // 默认公司
	}

	departments, err := h.service.GetTree(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, departments)
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var dept department.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&dept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dept)
}

func (h *DepartmentHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dept, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, dept)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// 使用 map 来检测哪些字段被显式提供（包括 null）
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建更新数据
	updates := map[string]interface{}{}
	if name, ok := req["name"].(string); ok && name != "" {
		updates["name"] = name
	}

	// 检查 parent_id 是否被显式提供（即使是 null）
	if _, ok := req["parent_id"]; ok {
		if parentID, ok := req["parent_id"].(float64); ok {
			updates["parent_id"] = int64(parentID)
		} else if req["parent_id"] == nil {
			updates["parent_id"] = nil
		}
	}

	// 检查 leader_id 是否被显式提供（即使是 null）
	if _, ok := req["leader_id"]; ok {
		if leaderID, ok := req["leader_id"].(float64); ok {
			updates["leader_id"] = int64(leaderID)
		} else if req["leader_id"] == nil {
			updates["leader_id"] = nil
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有需要更新的字段"})
		return
	}

	if err := h.service.UpdateFields(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的数据
	updated, _ := h.service.GetByID(id)
	c.JSON(http.StatusOK, updated)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// 获取部门信息，检查是否是最后一个
	dept, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "部门不存在"})
		return
	}

	// 检查是否是最后一个部门
	count, err := h.service.Count(dept.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法删除最后一个部门"})
		return
	}

	// 检查是否有子部门
	childIDs := h.service.GetAllChildIDs(id)
	hasChildren := len(childIDs) > 1

	// 获取转移目标部门
	var req struct {
		TransferToDeptID int64 `json:"transfer_to_dept_id"`
	}
	c.ShouldBindJSON(&req) // 忽略错误，空值用 0

	// 如果有子部门但没有指定转移目标，拒绝删除
	if hasChildren && req.TransferToDeptID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该部门下有子部门，请选择接收部门后再删除"})
		return
	}

	if err := h.service.Delete(id, req.TransferToDeptID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type ApprovalChainRequest struct {
	Steps []ApprovalStep `json:"steps"`
}

type ApprovalStep struct {
	EmployeeID int64 `json:"employee_id"`
	StepOrder int   `json:"step_order"`
}

func (h *DepartmentHandler) GetApprovalChain(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	chain, err := h.service.GetApprovalChain(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chain)
}

func (h *DepartmentHandler) SetApprovalChain(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req ApprovalChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chain := make([]department.DepartmentApprovalChain, len(req.Steps))
	for i, step := range req.Steps {
		chain[i] = department.DepartmentApprovalChain{
			DepartmentID: id,
			EmployeeID:   step.EmployeeID,
			StepOrder:    step.StepOrder,
		}
	}

	if err := h.service.SetApprovalChain(id, chain); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "approval chain updated"})
}
