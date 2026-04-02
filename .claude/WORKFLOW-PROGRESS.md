# Workflow 审批系统 - 项目进度

## 最近更新
2026-04-02

## 当前状态

### 已完成功能

| 功能 | 状态 | 文件更改 |
|-----|------|---------|
| 条件分支流转 | ✅ 完成 | `workflow_engine.go` - findNextNode 支持 yes/no 分支 |
| 审批人搜索 | ✅ 完成 | `employee.go` - SearchByName, `approval.go` - Search handler |
| 流程版本快照 | ✅ 完成 | `instance.go` - GraphData 字段保存运行时快照 |
| 撤回功能 | ✅ 完成 | `workflow.go` - CancelInstance 增加权限验证 |
| 会签节点 | ✅ 完成 | `CountersignNode.tsx` - 新增紫色会签节点组件 |
| 草稿自动保存 | ✅ 完成 | `WorkflowDesigner.tsx` - localStorage 自动保存 |
| 审批历史 | ✅ 完成 | `MyApplications.tsx` - Timeline 显示审批记录 |
| 批量审批 | ✅ 完成 | `approval.go` - BatchApprove/BatchReject |
| 发起流程入口 | ✅ 完成 | `DefinitionList/index.tsx` - 已发布流程显示"申请"按钮 |
| 数据权限隔离 | ✅ 完成 | `InstanceRepository` - GetPendingTasks/GetHandledTasks 增加 company_id 过滤 |
| 加签/减签 | ✅ 完成 | `ApprovalService` - AddApprover/RemoveApprover 方法 |
| 审批超时处理 | ✅ 完成 | `scheduler/scheduler.go` - TimeoutScheduler 定时检查超时的待审批任务 |

### 待完成功能

| 功能 | 优先级 | 说明 |
|-----|--------|------|
| 前端加签/减签 UI | P2 | 需要在 MyTasks 页面添加加签/减签按钮和对话框 |

## 已修复问题 (Review 2026-04-02)

| 问题 | 严重性 | 修复 |
|-----|--------|------|
| ApprovalService 未调用引擎推进流程 | 🔴 BLOCKER | Approve/Reject 改为调用 engine.ProcessApproval |
| 批量审批缺少权限验证 | 🔴 BLOCKER | BatchApprove/BatchReject 增加 approverID JWT 验证 |
| ListPending/ListHandled hardcoded user_id=1 | 🟡 WARNING | 改为从 JWT 获取当前用户ID |

## 待修复问题

1. **Docker 镜像未重建** - 网络问题导致构建失败，需要网络恢复后重建
   ```bash
   cd C:/project/Workflow_claude
   docker compose build backend frontend
   docker compose up -d
   ```

## 主要文件更改

### 后端
- `backend/internal/service/engine/workflow_engine.go` - 条件分支流转 + 超时截止时间设置
- `backend/internal/repository/employee.go` - SearchByName + GetCompanyIDByEmployeeID
- `backend/internal/repository/instance.go` - GetPendingTasks/GetHandledTasks 增加 company_id 过滤 + DeleteTask + GetOverduePendingTasks
- `backend/internal/service/approval.go` - AddApprover/RemoveApprover + ListPending/ListHandled 增加 companyID
- `backend/internal/service/scheduler/scheduler.go` - 新文件：超时调度器
- `backend/internal/domain/task/task.go` - DeadlineAt 字段
- `backend/internal/handler/api/v1/approval.go` - AddApprover/RemoveApprover handlers + company_id 过滤
- `backend/cmd/server/main.go` - 新增路由和调度器启动

### 前端
- `frontend/src/api/approval.ts` - addApprover/removeApprover API

## 环境信息

- Go 1.26 + Gin + GORM
- React 18 + Vite + Ant Design + @xyflow/react v12
- PostgreSQL + Redis + MinIO
- Docker Compose 部署
