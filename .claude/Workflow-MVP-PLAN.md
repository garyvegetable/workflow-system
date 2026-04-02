# Workflow MVP 实施计划

## 阶段划分

> ⚡ AI 加速：基于 gstack + Superpowers 框架，多子代理并行执行，大幅缩短工期。

| 阶段 | AI执行时间 | 内容 |
|------|-----------|------|
| **Phase 1** | 1-2 小时 | 项目初始化 + 数据库设计 |
| **Phase 2** | 1-2 小时 | 组织架构管理（公司/部门/员工/审批链） |
| **Phase 3** | 1-2 小时 | 基础数据管理（供应商/费用科目/银行账户） |
| **Phase 4** | 3-4 小时 | 流程设计器（React Flow + 节点配置） |
| **Phase 5** | 2-3 小时 | 动态表单 + 文件上传 |
| **Phase 6** | 2-3 小时 | 审批流程引擎 |
| **Phase 7** | 1-2 小时 | 通知系统 + 模板分发 |
| **Phase 8** | 1-2 小时 | 联调 + 测试 + 部署 |

**AI 连续工作总计：约 12-20 小时（1-2 天）**

---

## 子代理并行策略

### 什么时候用子代理

| 任务类型 | 适合子代理？ | 原因 |
|----------|-------------|------|
| 搜索/读取代码 | ✅ 是 | 用新鲜上下文，不污染主上下文 |
| 写重复性代码 | ✅ 是 | 并行加速 |
| 写单元测试 | ✅ 是 | 独立任务 |
| 架构设计/决策 | ❌ 否 | 需要主上下文的全部信息 |
| Plan 制定 | ❌ 否 | 需要综合判断 |
| Code Review | ❌ 否 | 需要完整上下文 |
| 集成测试 | ❌ 否 | 需要协调 |

### 典型并行模式

```
Phase 1 (并行):
├── 子代理 A: 初始化后端项目结构
├── 子代理 B: 初始化前端项目结构
└── 子代理 C: 编写 docker-compose.yml

Phase 2 (串行，需上下文):
└── 主代理: 依次实现公司→部门→员工 API + 前端
```

---

## 执行时间表

```
Hour 1-2:   并行子代理初始化 + 数据库设计
Hour 3-4:   子代理 CRUD API (公司/部门/员工)
Hour 5-6:   子代理 CRUD API (供应商/科目/银行)
Hour 7-10:  流程设计器 (主代理 + 子代理并行节点)
Hour 11-13: 动态表单 + 审批引擎
Hour 14-15: 通知 + 部署配置
Hour 16:    联调 + 验收
```

---

## Phase 1: 项目初始化

### 任务 1.1: 项目结构初始化

```
✅ 初始化后端项目
- 创建 backend/cmd/server/main.go
- 创建 backend/go.mod
- 创建 config/config.yaml
- 创建 internal/pkg/ 项目结构
- 验证: go build 成功

✅ 初始化前端项目
- 创建 frontend/src/ 结构
- 安装依赖: react, antd, react-flow, react-hook-form, axios
- 验证: npm run dev 成功

✅ Docker 配置
- 创建 Dockerfile (backend)
- 创建 Dockerfile (frontend)
- 创建 docker-compose.yml
- 创建 nginx.conf
- 验证: docker-compose up -d 成功
```

### 任务 1.2: 数据库设计

```
✅ 设计并创建所有表
- 创建 migrations/001_init/up.sql
- 包含所有表: company, department, employee, employee_department,
  department_approval_chain, supplier, employee_bank_account,
  expense_category, workflow_definition, workflow_instance,
  approval_task, attachment, notification, audit_log
- 创建索引
- 验证: psql 连接成功，表的创建

✅ 编写数据库迁移文档
```

---

## Phase 2: 组织架构管理

### 任务 2.1: 公司管理

```
✅ 后端 API
- internal/domain/company/company.go
- internal/repository/company/company.go
- internal/service/company/company.go
- internal/handler/api/v1/company.go
- API: CRUD /api/v1/companies

✅ 前端页面
- frontend/src/pages/company/CompanyList.tsx
- frontend/src/pages/company/CompanyForm.tsx
- frontend/src/api/company.ts

✅ 测试验证
- 可以创建/编辑/禁用公司
```

### 任务 2.2: 部门管理

```
✅ 后端 API
- internal/domain/department/department.go
- internal/repository/department/department.go
- internal/service/department/department.go
- internal/handler/api/v1/department.go
- API: CRUD /api/v1/departments

✅ 前端页面
- frontend/src/pages/department/DepartmentList.tsx (树形)
- frontend/src/pages/department/DepartmentForm.tsx

✅ 测试验证
- 可以创建部门（支持树形结构）
- 部门列表以树形展示
```

### 任务 2.3: 员工管理

```
✅ 后端 API
- internal/domain/employee/employee.go
- internal/repository/employee/employee.go
- internal/service/employee/employee.go
- internal/handler/api/v1/employee.go
- API: CRUD /api/v1/employees

✅ 前端页面
- frontend/src/pages/employee/EmployeeList.tsx
- frontend/src/pages/employee/EmployeeForm.tsx

✅ 测试验证
- 可以创建员工并分配多个部门
- 可以设置员工岗位级别
```

### 任务 2.4: 部门审批链

```
✅ 后端 API
- internal/repository/department_approval_chain/
- internal/service/department_approval_chain/
- API: GET/PUT /api/v1/departments/:id/approval-chain

✅ 前端
- 部门表单中增加审批链配置
- 可拖拽排序审批人

✅ 测试验证
- 可以为部门配置审批链顺序
```

---

## Phase 3: 基础数据管理

### 任务 3.1: 供应商管理

```
✅ 后端 API
- internal/domain/supplier/supplier.go
- internal/repository/supplier/supplier.go
- internal/service/supplier/supplier.go
- internal/handler/api/v1/supplier.go
- API: CRUD /api/v1/suppliers

✅ 前端页面
- frontend/src/pages/supplier/SupplierList.tsx
- frontend/src/pages/supplier/SupplierForm.tsx

✅ 测试验证
- 可以 CRUD 供应商
```

### 任务 3.2: 费用科目管理

```
✅ 后端 API
- internal/domain/expense_category/
- internal/repository/expense_category/
- internal/service/expense_category/
- API: CRUD /api/v1/expense-categories (树形)

✅ 前端页面
- frontend/src/pages/expense-category/ExpenseCategoryList.tsx
- frontend/src/pages/expense-category/ExpenseCategoryForm.tsx

✅ 测试验证
- 可以 CRUD 费用科目（支持层级）
```

### 任务 3.3: 员工银行账户

```
✅ 后端 API
- internal/domain/employee_bank_account/
- internal/repository/employee_bank_account/
- API: GET/POST/PUT/DELETE /api/v1/employees/:id/bank-accounts

✅ 前端
- 员工详情页增加银行账户 Tab

✅ 测试验证
- 可以为员工添加多个银行账户
- 可以设置默认账户
```

---

## Phase 4: 流程设计器

### 任务 4.1: React Flow 集成

```
✅ 安装依赖
- @xyflow/react (React Flow v12)

✅ 核心组件
- frontend/src/components/workflow/Designer/WorkflowCanvas.tsx
- frontend/src/components/workflow/Designer/NodePalette.tsx
- frontend/src/components/workflow/Designer/PropertyPanel.tsx

✅ 自定义节点
- StartNode.tsx
- ApprovalNode.tsx
- ConditionNode.tsx
- EndNode.tsx

✅ 测试验证
- 可以在画布上拖拽节点
- 可以连线
- 可以拖动位置
```

### 任务 4.2: 审批节点配置

```
✅ 审批节点属性面板
- 选择审批人类型（指定人员/部门审批链/角色）
- 选择审批人/部门/角色
- 配置通知角色

✅ 数据结构
- graph_data JSONB 存储

✅ 测试验证
- 可以配置审批节点
- 配置可以保存和加载
```

### 任务 4.3: 条件节点配置

```
✅ 条件节点属性面板
- 选择字段（从表单字段中选择）
- 选择运算符
- 输入比较值
- 配置 true/false 分支

✅ 条件表达式白名单验证
- internal/pkg/expression/validator.go
- 只允许表单字段名

✅ 测试验证
- 可以配置条件节点
- 条件可以正确解析
```

### 任务 4.4: 表单字段设计器

```
✅ 字段类型
- text, textarea, number, date, select, radio, checkbox, file, table

✅ 字段属性
- 名称、标签、必填、选项、默认值

✅ 明细表字段
- 可内嵌子字段

✅ 前端组件
- frontend/src/components/form/FieldEditor.tsx
- frontend/src/components/form/FieldRenderer.tsx

✅ 测试验证
- 可以设计各种类型的字段
- 可以预览表单效果
```

### 任务 4.5: 模板发布/版本

```
✅ 模板状态机
- 草稿 → 已发布 → 禁用

✅ 版本管理
- workflow_definition.version

✅ 前端
- 发布/禁用按钮
- 版本历史查看

✅ 测试验证
- 可以发布模板
- 草稿修改不影响已发布版本
```

---

## Phase 5: 动态表单 + 文件上传

### 任务 5.1: 动态表单渲染

```
✅ 表单渲染器
- frontend/src/components/form/DynamicForm.tsx
- 使用 react-hook-form
- 根据 form_fields JSON 动态渲染

✅ 表单验证
- 必填、格式、范围

✅ 明细表渲染
- 可动态添加/删除行

✅ 测试验证
- 选择模板后动态渲染表单
- 可以填写并提交
```

### 任务 5.2: 文件上传 (MinIO)

```
✅ 后端存储服务
- internal/pkg/storage/s3.go (统一接口)
- internal/pkg/storage/minio.go
- 支持切换: minio, s3, oss

✅ 文件上传 API
- POST /api/v1/attachments/upload
- 支持 multipart/form-data
- 返回 file_id

✅ 文件预览/下载
- GET /api/v1/attachments/:id/preview
- GET /api/v1/attachments/:id/download

✅ 限制
- 单文件 10MB
- 总附件 50MB

✅ 测试验证
- 可以上传 PDF/Excel/图片
- 可以预览和下载
```

---

## Phase 6: 审批流程引擎

### 任务 6.1: 发起流程

```
✅ 流程发起服务
- internal/service/engine/workflow_engine.go
- internal/service/engine/assignee_resolver.go

✅ 审批人解析
- 根据 department_approval_chain 查询下一个审批人
- 驳回后重提从头开始

✅ 创建实例
- workflow_instance 表

✅ 创建任务
- approval_task 表

✅ 测试验证
- 提交表单后生成审批任务
- 审批人可以在待办列表中看到
```

### 任务 6.2: 审批动作

```
✅ 同意
- internal/service/engine/workflow_engine.go Approve()
- 更新任务状态
- 流转到下一节点或结束

✅ 驳回
- internal/service/engine/workflow_engine.go Reject()
- 更新实例状态为草稿

✅ 转签
- internal/service/engine/workflow_engine.go Transfer()
- 更新任务审批人

✅ 并发控制
- 乐观锁 (version 检查)

✅ 测试验证
- 可以同意/驳回
- 驳回后可以修改重提
- 并发审批被阻止
```

### 任务 6.3: 条件分支处理

```
✅ 条件表达式执行
- 使用 govaluate
- 白名单验证字段名

✅ 分支流转
- 根据条件结果选择下一节点

✅ 测试验证
- 金额 > 20万走高层审批
- 金额 <= 20万走普通审批
```

### 任务 6.4: 前端审批页面

```
✅ 待办任务列表
- frontend/src/pages/approval/MyTasks.tsx

✅ 审批操作
- frontend/src/components/workflow/TaskAction/TaskAction.tsx
- 同意/驳回按钮
- 审批意见输入

✅ 审批历史
- frontend/src/pages/approval/InstanceHistory.tsx

✅ 测试验证
- 审批人可以看到待办任务
- 可以执行审批动作
- 可以查看审批历史
```

---

## Phase 7: 通知系统 + 模板分发

### 任务 7.1: 邮件通知 (Mock)

```
✅ 邮件服务
- internal/pkg/email/email.go (Mock)
- 后续可切换到真实 SMTP

✅ 通知触发点
- 新申请时通知审批人
- 审批结果时通知申请人
- 驳回时通知申请人

✅ 测试验证
- 控制台输出邮件内容 (Mock)
```

### 任务 7.2: 站内通知

```
✅ 通知表
- notification 表

✅ API
- GET /api/v1/notifications
- PUT /api/v1/notifications/:id/read
- PUT /api/v1/notifications/read-all

✅ 前端
- frontend/src/pages/notification/NotificationList.tsx
- 通知中心入口

✅ 测试验证
- 可以看到站内通知
```

### 任务 7.3: 模板分发

```
✅ 复制 API
- POST /api/v1/workflows/:id/copy
- 复制到指定公司

✅ 前端
- 模板列表增加"复制"按钮

✅ 测试验证
- 可以将模板复制到其他公司
```

---

## Phase 8: 联调 + 测试 + 部署

### 任务 8.1: 前后端联调

```
✅ 认证流程
- JWT 登录/登出
- 接口权限校验

✅ 完整流程测试
- 创建公司 → 创建部门 → 创建员工 → 配置审批链
- 创建供应商/费用科目
- 设计流程模板 → 发布
- 提交申请 → 审批 → 完成
```

### 任务 8.2: 安全加固

```
✅ SQL 注入防护
- 参数化查询 (GORM)

✅ XSS 防护
- React 自动转义

✅ CORS 配置
- nginx cors

✅ 文件上传安全
- 白名单 MIME 类型
```

### 任务 8.3: 部署

```
✅ Docker 镜像构建
- backend Dockerfile
- frontend Dockerfile

✅ docker-compose.yml 完善
- 环境变量配置
- 健康检查

✅ 生产部署文档
- README.md 更新
```

---

## 验收标准清单

### 组织管理
- [ ] 可以创建/编辑/禁用公司
- [ ] 可以创建部门（树形结构）
- [ ] 可以创建员工并分配部门和级别
- [ ] 一个员工可以属于多个部门
- [ ] 可以管理供应商（CRUD）
- [ ] 可以管理费用科目（CRUD、层级结构）
- [ ] 可以记录员工银行账户

### 流程设计
- [ ] 可视化拖拽设计流程图
- [ ] 可以添加审批节点并指定审批人/部门/级别
- [ ] 可以添加条件节点（金额判断）
- [ ] 可以定义表单字段（含明细表）
- [ ] 可以发布/禁用模板

### 申请提交
- [ ] 选择模板后动态渲染表单
- [ ] 可以上传附件（PDF/Excel/图片）
- [ ] 提交后生成审批任务

### 审批
- [ ] 审批人可以看到待办任务
- [ ] 可以同意/驳回
- [ ] 驳回后可以修改重提
- [ ] 可以转签给其他人
- [ ] 可以查看审批历史

### 通知
- [ ] 新申请时审批人收到邮件 (Mock)
- [ ] 审批结果时申请人收到邮件 (Mock)
- [ ] 站内通知列表

### 模板分发
- [ ] 可以将模板复制到其他公司

---

## 里程碑

| 里程碑 | 日期 | 交付物 |
|--------|------|---------|
| M1 | 第 2 周末 | 组织架构管理完成 |
| M2 | 第 5 周末 | 流程设计器完成 |
| M3 | 第 7 周末 | 审批流程完成 |
| M4 | 第 9 周末 | MVP 交付 |

---

## 风险与缓解

| 风险 | 概率 | 影响 | 缓解 |
|------|------|------|------|
| React Flow 定制复杂 | 中 | 中 | 提前研究文档，先做简单节点 |
| 动态表单验证复杂 | 中 | 中 | 使用成熟的 react-hook-form |
| 条件表达式边界情况 | 低 | 中 | 充分单元测试 |
| 审批链配置用户不会用 | 中 | 中 | 提供默认配置引导 |
