package engine

import (
	"encoding/json"
	"fmt"
	"time"
	"workflow-system/internal/domain/instance"
	"workflow-system/internal/domain/task"
	"workflow-system/internal/repository"
)

var timeNow = time.Now

// WorkflowEngine 处理流程执行逻辑
type WorkflowEngine struct {
	workflowRepo  *repository.WorkflowRepository
	instanceRepo  *repository.InstanceRepository
	departmentRepo *repository.DepartmentRepository
	notifService  NotificationService
}

// NotificationService 通知服务接口
type NotificationService interface {
	SendToUser(userID int64, message string) error
}

// WorkflowEngineDeps 工作流引擎依赖
type WorkflowEngineDeps struct {
	WorkflowRepo  *repository.WorkflowRepository
	InstanceRepo  *repository.InstanceRepository
	DepartmentRepo *repository.DepartmentRepository
	NotifService  NotificationService
}

func NewWorkflowEngine(deps WorkflowEngineDeps) *WorkflowEngine {
	return &WorkflowEngine{
		workflowRepo:  deps.WorkflowRepo,
		instanceRepo:  deps.InstanceRepo,
		departmentRepo: deps.DepartmentRepo,
		notifService:  deps.NotifService,
	}
}

// StartWorkflow 启动流程实例
func (e *WorkflowEngine) StartWorkflow(defID int64, initiatorID int64, formData map[string]interface{}) (*instance.WorkflowInstance, error) {
	// 1. 获取流程定义
	def, err := e.workflowRepo.GetByID(defID)
	if err != nil {
		return nil, fmt.Errorf("流程定义不存在: %w", err)
	}
	if def.Status != 2 { // 2 = 已发布
		return nil, fmt.Errorf("流程未发布")
	}

	// 2. 解析图形数据
	var graphData map[string]interface{}
	if err := json.Unmarshal(def.GraphData, &graphData); err != nil {
		return nil, fmt.Errorf("流程图形数据解析失败: %w", err)
	}

	if graphData == nil {
		return nil, fmt.Errorf("流程图形数据无效")
	}

	nodes := graphData["nodes"].([]interface{})
	edges := graphData["edges"].([]interface{})

	// 3. 找到开始节点
	var startNode map[string]interface{}
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if node["type"] == "start" {
			startNode = node
			break
		}
	}
	if startNode == nil {
		return nil, fmt.Errorf("未找到开始节点")
	}

	// 4. 创建流程实例（保存流程定义快照，确保运行时使用版本不变）
	formDataJSON, _ := json.Marshal(formData)
	inst := &instance.WorkflowInstance{
		DefinitionID: defID,
		Title:        def.Name,
		Status:       1, // 审批中
		InitiatorID:  initiatorID,
		CompanyID:    def.CompanyID,
		FormData:     formDataJSON,
		GraphData:    def.GraphData, // 保存流程定义快照
	}
	if err := e.instanceRepo.Create(inst); err != nil {
		return nil, err
	}

	// 5. 创建初始任务（通常是第一个审批节点）
	// 查找开始节点后的第一个审批节点（从"是"分支开始）
	nextNode := e.findNextNode(startNode, nodes, edges, "yes")
	if nextNode != nil {
		t := e.createTask(inst, nextNode, 0)
		if err := e.instanceRepo.CreateTask(t); err != nil {
			return nil, err
		}
		// 发送通知
		if e.notifService != nil && t.AssigneeID > 0 {
			e.notifService.SendToUser(t.AssigneeID, fmt.Sprintf("您有新的待处理任务: %s", inst.Title))
		}
	}

	return inst, nil
}

// ProcessApproval 处理审批结果
func (e *WorkflowEngine) ProcessApproval(taskID int64, approverID int64, result string, comment string) error {
	// 1. 获取任务
	t, err := e.instanceRepo.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("任务不存在: %w", err)
	}
	if t.Status != 1 { // 1 = 待处理
		return fmt.Errorf("任务已处理")
	}

	// 2. 更新任务状态
	t.Status = 2 // 已处理
	t.Action = result
	t.Comment = comment
	now := timeNow().Unix()
	t.CompletedAt = &now
	if err := e.instanceRepo.UpdateTask(t); err != nil {
		return err
	}

	// 3. 获取流程实例（使用实例中保存的流程定义快照）
	inst, err := e.instanceRepo.GetByID(t.InstanceID)
	if err != nil {
		return err
	}

	// 4. 根据审批结果处理（使用实例的流程定义快照，而非最新定义）
	var graphData map[string]interface{}
	if err := json.Unmarshal(inst.GraphData, &graphData); err != nil {
		return fmt.Errorf("流程图形数据解析失败: %w", err)
	}

	if graphData == nil {
		return fmt.Errorf("流程图形数据无效")
	}

	nodes := graphData["nodes"].([]interface{})
	edges := graphData["edges"].([]interface{})

	currentNode := e.findNodeByID(t.NodeID, nodes)

	if result == "approve" {
		// 同意 - 查找下一个节点（根据条件分支）
		nextNode := e.findNextNode(currentNode, nodes, edges, "yes")
		if nextNode != nil {
			nodeType, _ := nextNode["type"].(string)
			if nodeType == "end" {
				// 流程结束
				inst.Status = 2 // 已通过
				e.instanceRepo.Update(inst)
				if e.notifService != nil {
					e.notifService.SendToUser(inst.InitiatorID, fmt.Sprintf("您的申请已通过: %s", inst.Title))
				}
			} else {
				// 创建下一个任务
				newTask := e.createTask(inst, nextNode, 0)
				e.instanceRepo.CreateTask(newTask)
				if e.notifService != nil && newTask.AssigneeID > 0 {
					e.notifService.SendToUser(newTask.AssigneeID, fmt.Sprintf("您有新的待处理任务: %s", inst.Title))
				}
			}
		}
	} else {
		// 驳回 - 查找条件节点的"否"分支
		nextNode := e.findNextNode(currentNode, nodes, edges, "no")
		if nextNode != nil {
			nodeType, _ := nextNode["type"].(string)
			if nodeType == "end" {
				// 驳回也结束流程
				inst.Status = 3 // 已驳回
				e.instanceRepo.Update(inst)
				if e.notifService != nil {
					e.notifService.SendToUser(inst.InitiatorID, fmt.Sprintf("您的申请已被驳回: %s", inst.Title))
				}
			} else {
				// 创建下一个任务（驳回分支）
				newTask := e.createTask(inst, nextNode, 0)
				e.instanceRepo.CreateTask(newTask)
				if e.notifService != nil && newTask.AssigneeID > 0 {
					e.notifService.SendToUser(newTask.AssigneeID, fmt.Sprintf("您有新的待处理任务（驳回分支）: %s", inst.Title))
				}
			}
		} else {
			// 没有驳回分支，直接结束
			inst.Status = 3 // 已驳回
			e.instanceRepo.Update(inst)
			if e.notifService != nil {
				e.notifService.SendToUser(inst.InitiatorID, fmt.Sprintf("您的申请已被驳回: %s", inst.Title))
			}
		}
	}

	return nil
}

// findNextNode 查找下一个节点
// branch: "yes" 或 "no"，用于条件节点选择分支
func (e *WorkflowEngine) findNextNode(current map[string]interface{}, nodes []interface{}, edges []interface{}, branch string) map[string]interface{} {
	currentID, _ := current["id"].(string)
	nodeType, _ := current["type"].(string)

	// 如果是条件节点，优先使用 branch 参数查找对应分支
	if nodeType == "condition" && edges != nil && len(edges) > 0 {
		for _, edgeRaw := range edges {
			edge := edgeRaw.(map[string]interface{})
			source, _ := edge["source"].(string)
			sourceHandle, _ := edge["sourceHandle"].(string)
			if source == currentID && sourceHandle == branch {
				target, _ := edge["target"].(string)
				return e.findNodeByID(target, nodes)
			}
		}
		// 如果没找到对应分支的边，尝试找默认分支
		if branch == "yes" {
			return e.findNextNodeByDefault(current, nodes, edges)
		}
		return nil
	}

	// 非条件节点或没有边的情况，使用边查找
	if edges != nil && len(edges) > 0 {
		for _, edgeRaw := range edges {
			edge := edgeRaw.(map[string]interface{})
			source, _ := edge["source"].(string)
			if source == currentID {
				target, _ := edge["target"].(string)
				return e.findNodeByID(target, nodes)
			}
		}
	}

	// 降级：使用默认顺序查找
	return e.findNextNodeByDefault(current, nodes, edges)
}

// findNextNodeByDefault 使用默认顺序查找下一个节点
func (e *WorkflowEngine) findNextNodeByDefault(current map[string]interface{}, nodes []interface{}, edges []interface{}) map[string]interface{} {
	currentFound := false
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if currentFound {
			nodeType, _ := node["type"].(string)
			if nodeType == "approval" || nodeType == "condition" {
				return node
			}
			if nodeType == "end" {
				return node
			}
		}
		if node["id"] == current["id"] {
			currentFound = true
		}
	}
	return nil
}

// findNodeByID 根据ID查找节点
func (e *WorkflowEngine) findNodeByID(nodeID string, nodes []interface{}) map[string]interface{} {
	for _, n := range nodes {
		node := n.(map[string]interface{})
		if node["id"] == nodeID {
			return node
		}
	}
	return nil
}

// createTask 创建任务
func (e *WorkflowEngine) createTask(inst *instance.WorkflowInstance, node map[string]interface{}, approverID int64) *task.ApprovalTask {
	nodeData, _ := node["data"].(map[string]interface{})
	nodeName := ""
	if nodeData != nil {
		if name, ok := nodeData["label"].(string); ok {
			nodeName = name
		} else if name, ok := nodeData["name"].(string); ok {
			nodeName = name
		}
	}
	if nodeName == "" {
		nodeName, _ = node["id"].(string)
	}

	t := &task.ApprovalTask{
		InstanceID: inst.ID,
		NodeID:     node["id"].(string),
		NodeName:   nodeName,
		Status:     1, // 待处理
	}

	// 如果没有指定审批人，从节点配置获取
	if approverID == 0 {
		if nodeData != nil {
			// 优先使用直接指定的审批人
			if assigneeID, ok := nodeData["assigneeId"].(float64); ok && assigneeID > 0 {
				t.AssigneeID = int64(assigneeID)
			} else if deptID, ok := nodeData["departmentId"].(float64); ok && deptID > 0 {
				// 根据部门ID获取部门负责人
				if leaderID, err := e.resolveDepartmentLeader(int64(deptID)); err == nil && leaderID > 0 {
					t.AssigneeID = leaderID
				}
			} else if assigneeName, ok := nodeData["approverName"].(string); ok && assigneeName != "" {
				// TODO: 根据姓名查找用户ID
				// 暂时留空，让后续处理人指定
			}
		}
	} else {
		t.AssigneeID = approverID
	}

	// 设置超时时间（节点配置中的 deadline_hours，单位小时）
	if nodeData != nil {
		if deadlineHours, ok := nodeData["deadlineHours"].(float64); ok && deadlineHours > 0 {
			deadline := timeNow().Add(time.Duration(deadlineHours) * time.Hour).Unix()
			t.DeadlineAt = &deadline
		}
	}

	return t
}

// resolveDepartmentLeader 根据部门ID获取部门负责人
func (e *WorkflowEngine) resolveDepartmentLeader(deptID int64) (int64, error) {
	dept, err := e.departmentRepo.GetByID(deptID)
	if err != nil {
		return 0, fmt.Errorf("部门不存在: %w", err)
	}
	if dept.LeaderID == nil || *dept.LeaderID == 0 {
		return 0, fmt.Errorf("部门未指定负责人")
	}
	return *dept.LeaderID, nil
}
