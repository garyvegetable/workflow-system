# Workflow 系统 MVP 设计文档

## 背景

现有系统 norning 极度繁琐、维护困难、成本高昂。需求一套自研的企业内部审批流程管理系统，具备高度自定义的流程模板能力。

## 核心痛点

1. **norning 维护困难** — 流程配置复杂，管理员操作门槛高
2. **审批链路不透明** — 员工不清楚流程走到哪一步
3. **费用高** — 商业软件 license 费用昂贵
4. **不够灵活** — 无法快速响应业务变化

## 目标

构建一套**高度可自定义**的审批流程管理系统，让管理员可以通过可视化界面定义任意审批链路，员工可以便捷提交申请并实时追踪状态。

## 方案选择

**选择方案 C — 完整 MVP**

理由：需要完整功能才能说服老板，残缺功能会导致不信任。

## 功能列表

### 1. 组织架构管理

| 功能 | 说明 |
|------|------|
| 公司管理 | CRUD、状态启用/禁用 |
| 部门管理 | 树形结构、CRUD |
| 人员管理 | CRUD、关联多部门、岗位级别 |
| 岗位级别 | 创建员工时定义（普通员工/主管/经理/总监） |

### 2. 基础数据管理

| 功能 | 说明 |
|------|------|
| 供应商管理 | CRUD、状态启用/禁用（采购必需） |
| 费用科目管理 | CRUD、层级结构（办公用品/设备采购/差旅费等） |
| 员工银行账户 | 银行名称、账号、开户行（报销付款必需） |

### 3. 流程模板管理

| 功能 | 说明 |
|------|------|
| 可视化流程设计器 | 拖拽节点、连线、节点属性配置 |
| 节点类型 | 开始节点、审批节点、条件节点、结束节点 |
| 审批人配置 | 指定人员、**部门审批链（管理员配置）**、岗位级别 |
| 条件分支 | 金额 > X 时走向不同审批链 |
| 表单字段类型 | 文本、数字、日期、下拉、单选、多选、文件上传、**明细表（可内嵌多行）** |
| 表单字段 | 文本、数字、日期、文件上传 |
| 模板发布/禁用 | 版本管理 |

### 3.1 部门审批链

管理员在部门中配置审批顺序链：

```
部门 A 审批链：
1. 张三（员工）
2. 李四（主管）
3. 王五（经理）
```

- 审批时按顺序流转
- 驳回则终止，重提从头开始
- 部门负责人不自动设为审批人，由管理员指定

### 4. 申请提交

| 功能 | 说明 |
|------|------|
| 申请入口 | 选择模板 → 填写表单 → 提交 |
| 表单渲染 | 根据模板动态渲染字段 |
| 文件上传 | 支持 PDF、JPG、Excel、CSV |
| 文件预览/下载 | 审批人可以查看附件 |
| 提交历史 | 查看自己发起的申请 |

### 5. 审批流程

| 功能 | 说明 |
|------|------|
| 待办任务 | 列出需要审批的任务 |
| 审批动作 | 同意、驳回 |
| 驳回重提 | 修改后重新提交，已审批节点可跳过 |
| 审批历史 | 查看整个审批链路 |
| 转签 | 当前审批人转给其他人 |

### 6. 通知系统

| 功能 | 说明 |
|------|------|
| 邮件通知 | 新申请通知、审批结果通知、待办提醒 |
| 站内通知 | WebSocket 实时推送 |
| 通知中心 | 查看所有通知、已读/未读 |

### 7. 系统管理

| 功能 | 说明 |
|------|------|
| 模板复制 | 将模板复制到其他公司 |
| 审计日志 | 记录关键操作 |
| LDAP 集成 | 微软账号登录（Mock） |
| SMTP 配置 | 邮件发送配置 |

## 数据模型

### 公司 (company)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| code | VARCHAR(50) | 公司代码 |
| name | VARCHAR(200) | 公司名称 |
| short_name | VARCHAR(100) | 简称 |
| status | SMALLINT | 1:正常 2:禁用 |
| schema_name | VARCHAR(100) | 数据库 schema |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 部门 (department)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 所属公司 |
| parent_id | BIGINT | 上级部门 |
| name | VARCHAR(100) | 部门名称 |
| leader_id | BIGINT | 负责人 |
| status | SMALLINT | 1:正常 2:禁用 |
| created_at | TIMESTAMP | 创建时间 |

### 部门审批链 (department_approval_chain)

管理员为部门配置审批顺序链。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| department_id | BIGINT | 部门ID |
| employee_id | BIGINT | 审批人ID |
| step_order | INT | 审批顺序（1, 2, 3...） |
| created_at | TIMESTAMP | 创建时间 |

### 员工 (employee)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 所属公司 |
| username | VARCHAR(100) | 用户名 |
| email | VARCHAR(255) | 邮箱 |
| password_hash | VARCHAR(255) | 密码 |
| level | VARCHAR(50) | 岗位级别 |
| status | SMALLINT | 1:正常 2:禁用 |
| created_at | TIMESTAMP | 创建时间 |

### 员工部门关联 (employee_department)

| 字段 | 类型 | 说明 |
|------|------|------|
| employee_id | BIGINT | 员工ID |
| department_id | BIGINT | 部门ID |

### 供应商 (supplier)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 所属公司 |
| code | VARCHAR(50) | 供应商代码 |
| name | VARCHAR(200) | 供应商名称 |
| contact | VARCHAR(100) | 联系人 |
| phone | VARCHAR(50) | 电话 |
| email | VARCHAR(255) | 邮箱 |
| address | VARCHAR(500) | 地址 |
| bank_name | VARCHAR(200) | 开户银行 |
| bank_account | VARCHAR(100) | 银行账号 |
| tax_number | VARCHAR(50) | 税号 |
| status | SMALLINT | 1:正常 2:禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 员工银行账户 (employee_bank_account)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| employee_id | BIGINT | 员工ID |
| bank_name | VARCHAR(200) | 开户银行 |
| bank_branch | VARCHAR(200) | 支行名称 |
| bank_account | VARCHAR(100) | 银行账号 |
| account_holder | VARCHAR(100) | 开户人姓名 |
| is_default | BOOLEAN | 是否默认账户 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 费用科目 (expense_category)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 所属公司 |
| code | VARCHAR(50) | 科目代码 |
| name | VARCHAR(200) | 科目名称 |
| parent_id | BIGINT | 上级科目（一级留空） |
| status | SMALLINT | 1:正常 2:禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 流程定义 (workflow_definition)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 所属公司 |
| code | VARCHAR(50) | 模板代码 |
| name | VARCHAR(200) | 模板名称 |
| version | INT | 版本号 |
| graph_data | JSONB | 流程图数据 |
| form_fields | JSONB | 表单字段定义 |
| status | SMALLINT | 1:草稿 2:已发布 3:禁用 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 流程实例 (workflow_instance)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| definition_id | BIGINT | 流程定义ID |
| company_id | BIGINT | 所属公司 |
| title | VARCHAR(500) | 实例标题 |
| initiator_id | BIGINT | 发起人 |
| form_data | JSONB | 表单数据 |
| status | SMALLINT | 0:草稿 1:审批中 2:已通过 3:已驳回 4:已撤回 |
| current_nodes | JSONB | 当前节点 |
| started_at | TIMESTAMP | 开始时间 |
| finished_at | TIMESTAMP | 结束时间 |

### 审批任务 (approval_task)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| instance_id | BIGINT | 流程实例ID |
| node_id | VARCHAR(50) | 节点ID |
| node_name | VARCHAR(100) | 节点名称 |
| assignee_id | BIGINT | 审批人 |
| status | SMALLINT | 1:待审批 2:已审批 3:已驳回 |
| action | VARCHAR(20) | approve/reject |
| comment | TEXT | 审批意见 |
| completed_at | TIMESTAMP | 完成时间 |
| version | INT | 乐观锁版本号（防止并发审批） |
| created_at | TIMESTAMP | 创建时间 |

### 附件 (attachment)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| instance_id | BIGINT | 流程实例ID |
| field_name | VARCHAR(100) | 字段名 |
| file_name | VARCHAR(255) | 文件名 |
| file_path | VARCHAR(500) | 文件路径 |
| file_size | BIGINT | 文件大小 |
| mime_type | VARCHAR(100) | MIME类型 |
| uploaded_at | TIMESTAMP | 上传时间 |

### 通知 (notification)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| user_id | BIGINT | 接收人 |
| title | VARCHAR(200) | 通知标题 |
| content | TEXT | 通知内容 |
| type | VARCHAR(50) | 类型 |
| is_read | BOOLEAN | 已读 |
| created_at | TIMESTAMP | 创建时间 |

### 审计日志 (audit_log)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| company_id | BIGINT | 公司ID |
| user_id | BIGINT | 用户ID |
| action | VARCHAR(100) | 操作 |
| resource_type | VARCHAR(100) | 资源类型 |
| resource_id | VARCHAR(100) | 资源ID |
| details | JSONB | 详情 |
| ip_address | VARCHAR(50) | IP |
| created_at | TIMESTAMP | 时间 |

## 技术决策

### 多租户隔离策略

**采用逻辑隔离（company_id 字段）**

- 所有表通过 `company_id` 字段隔离数据
- API 层强制校验 `company_id`，禁止跨公司访问
- 适合 MVP 阶段，生产环境如需高安全可升级为 Schema 隔离

### 审批人解析

**采用管理员配置的审批链（department_approval_chain）**

- 管理员在部门中配置审批顺序链
- 审批时按顺序流转，已审批人不会重复审批
- 驳回则终止审批链，重提从头开始

### 条件表达式安全

**白名单字段验证**

- 条件表达式只允许访问表单中定义的字段名
- 非法字段名（如 `env.*`、`file.*`）直接返回 false
- 支持运算符：`>`, `<`, `==`, `!=`, `>=`, `<=`

### 并发控制

**乐观锁（version 字段）**

- `approval_task` 表含 `version` 字段
- 审批时检查 version，如已变更则拒绝操作
- 前端提示"该任务已被其他用户处理"

### 文件上传限制

| 限制项 | 值 |
|--------|-----|
| 单文件大小 | 10MB |
| 总附件大小 | 50MB |
| 支持格式 | PDF, JPG, PNG, GIF, Excel, CSV |

### 驳回重提状态机

```
驳回后：
1. 实例状态改为"草稿"（status=0）
2. 已审批任务保留历史记录
3. 发起人可修改表单重新提交
4. 重新提交后，新审批链从头开始
```

## API 设计

### 认证

```
POST   /api/v1/auth/login          # 登录
POST   /api/v1/auth/logout         # 登出
GET    /api/v1/auth/current        # 获取当前用户
```

### 组织管理

```
GET    /api/v1/companies                      # 公司列表
POST   /api/v1/companies                     # 创建公司
GET    /api/v1/companies/:id                 # 公司详情
PUT    /api/v1/companies/:id                 # 更新公司

GET    /api/v1/departments                   # 部门列表（树形）
POST   /api/v1/departments                   # 创建部门
GET    /api/v1/departments/:id               # 部门详情
PUT    /api/v1/departments/:id               # 更新部门

GET    /api/v1/departments/:id/approval-chain # 获取部门审批链
PUT    /api/v1/departments/:id/approval-chain # 设置部门审批链（覆盖）

GET    /api/v1/employees                     # 员工列表
POST   /api/v1/employees                     # 创建员工
GET    /api/v1/employees/:id                # 员工详情
PUT    /api/v1/employees/:id                # 更新员工
DELETE /api/v1/employees/:id                # 删除员工

GET    /api/v1/employees/:id/bank-accounts # 员工银行账户列表
POST   /api/v1/employees/:id/bank-accounts # 添加银行账户
PUT    /api/v1/employees/:id/bank-accounts/:aid # 更新银行账户
DELETE /api/v1/employees/:id/bank-accounts/:aid # 删除银行账户

### 供应商管理

```
GET    /api/v1/suppliers                     # 供应商列表
POST   /api/v1/suppliers                     # 创建供应商
GET    /api/v1/suppliers/:id                # 供应商详情
PUT    /api/v1/suppliers/:id                # 更新供应商
DELETE /api/v1/suppliers/:id                # 删除供应商

### 费用科目管理

```
GET    /api/v1/expense-categories             # 费用科目列表（树形）
POST   /api/v1/expense-categories             # 创建费用科目
GET    /api/v1/expense-categories/:id        # 费用科目详情
PUT    /api/v1/expense-categories/:id        # 更新费用科目
DELETE /api/v1/expense-categories/:id        # 删除费用科目
```

### 流程管理

```
GET    /api/v1/workflows                      # 模板列表
POST   /api/v1/workflows                      # 创建模板
GET    /api/v1/workflows/:id                 # 模板详情
PUT    /api/v1/workflows/:id                 # 更新模板
POST   /api/v1/workflows/:id/publish         # 发布模板
POST   /api/v1/workflows/:id/copy           # 复制到其他公司
DELETE /api/v1/workflows/:id                 # 删除模板

POST   /api/v1/workflows/instances           # 发起流程
GET    /api/v1/workflows/instances/:id      # 流程详情
POST   /api/v1/workflows/instances/:id/cancel # 撤回
```

### 审批

```
GET    /api/v1/tasks/pending                  # 待审批任务
GET    /api/v1/tasks/handled                 # 已审批任务
POST   /api/v1/tasks/:id/approve             # 同意
POST   /api/v1/tasks/:id/reject              # 驳回
POST   /api/v1/tasks/:id/transfer            # 转签
GET    /api/v1/tasks/:id/history             # 审批历史
```

### 附件

```
POST   /api/v1/attachments/upload            # 上传附件
GET    /api/v1/attachments/:id/download      # 下载附件
GET    /api/v1/attachments/:id/preview      # 预览附件
```

### 通知

```
GET    /api/v1/notifications                  # 通知列表
PUT    /api/v1/notifications/:id/read       # 标记已读
PUT    /api/v1/notifications/read-all       # 全部已读
```

### 系统

```
GET    /api/v1/audit-logs                    # 审计日志
GET    /api/v1/system/config                  # 系统配置
PUT    /api/v1/system/config                  # 更新配置
```

## 验收标准

### 组织管理

- [ ] 可以创建/编辑/禁用公司
- [ ] 可以创建部门（树形结构）
- [ ] 可以创建员工并分配部门和级别
- [ ] 一个员工可以属于多个部门
- [ ] 可以管理供应商（CRUD）
- [ ] 可以管理费用科目（CRUD、层级结构）
- [ ] 可以记录员工银行账户（报销付款用）

### 流程设计

- [ ] 可视化拖拽设计流程图
- [ ] 可以添加审批节点并指定审批人/部门/级别
- [ ] 可以添加条件节点（金额判断）
- [ ] 可以定义表单字段（含明细表/多行输入）
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

- [ ] 新申请时审批人收到邮件
- [ ] 审批结果时申请人收到邮件
- [ ] 站内通知实时推送

### 模板分发

- [ ] 可以将模板复制到其他公司

## 技术栈

### 后端

- Go 1.26 + Gin
- GORM + PostgreSQL
- JWT 认证
- MinIO (文件存储)

### 前端

- React 18 + TypeScript
- Vite + Ant Design 5.x
- Redux Toolkit
- React Flow (流程设计器)
- React Hook Form (表单)
- Axios

### 部署

- Docker + Docker Compose
- Nginx

## 工作量估算

| 模块 | 估算 |
|------|------|
| 组织架构管理 | 1 周 |
| 流程设计器 | 2 周 |
| 动态表单+附件 | 1 周 |
| 审批流程 | 1.5 周 |
| 通知系统 | 0.5 周 |
| 模板分发 | 0.5 周 |
| 联调+测试 | 1.5 周 |
| **总计** | **8 周** |
