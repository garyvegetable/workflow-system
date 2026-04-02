# Workflow MVP 技术研究

## 研究目标

为 workflow-system MVP 选择最佳技术方案，重点研究：
1. 流程设计器（React Flow vs 自研）
2. 动态表单方案（表单生成器）
3. 文件存储方案
4. 条件表达式引擎

---

## 1. 流程设计器

### 方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **React Flow** | 功能完整、文档好、社区活跃 | 包体积较大 (~200KB) | ⭐⭐⭐⭐ |
| **自研 Canvas** | 完全可控、定制灵活 | 开发周期长、bug 多 | ⭐⭐ |
| **GoJS** | 强大、类型丰富 | 商业授权、复杂 | ⭐⭐⭐ |
| **AntV X6** | 阿里系、轻量 | 中文文档少 | ⭐⭐⭐ |

### React Flow 选型理由

```
1. 开源免费 (MIT)
2. React 原生集成
3. 节点可高度自定义
4. 支持连线、拖拽、缩放
5. 社区成熟，问题易解决
6.与我们技术栈(React+TS)一致
```

### React Flow 核心概念

```typescript
// 节点类型
type NodeType = 'start' | 'approval' | 'condition' | 'end';

// 审批节点配置
interface ApprovalNode {
  type: 'approval';
  data: {
    assigneeType: 'user' | 'department_leader' | 'role';
    assigneeValue: string | number;
    notifyRoles?: string[];
  };
}

// 条件节点配置
interface ConditionNode {
  type: 'condition';
  data: {
    field: string;      // 字段名
    operator: '>' | '<' | '==' | '!=';
    value: any;         // 比较值
    trueTarget: string; // 条件为真时的下一节点
    falseTarget: string; // 条件为假时的下一节点
  };
}
```

### 流程图数据结构设计

```json
{
  "nodes": [
    { "id": "start_1", "type": "start", "position": { "x": 100, "y": 200 }, "data": {} },
    { "id": "approval_1", "type": "approval", "position": { "x": 300, "y": 200 }, "data": { "assigneeType": "department_leader" } },
    { "id": "condition_1", "type": "condition", "position": { "x": 500, "y": 200 }, "data": { "field": "amount", "operator": ">", "value": 200000 } },
    { "id": "approval_2", "type": "approval", "position": { "x": 700, "y": 100 }, "data": { "assigneeType": "role", "assigneeValue": "senior_manager" } },
    { "id": "approval_3", "type": "approval", "position": { "x": 700, "y": 300 }, "data": { "assigneeType": "department_leader" } },
    { "id": "end_1", "type": "end", "position": { "x": 900, "y": 200 }, "data": {} }
  ],
  "edges": [
    { "source": "start_1", "target": "approval_1" },
    { "source": "approval_1", "target": "condition_1" },
    { "source": "condition_1", "target": "approval_2", "label": "金额 > 20万" },
    { "source": "condition_1", "target": "approval_3", "label": "金额 <= 20万" },
    { "source": "approval_2", "target": "end_1" },
    { "source": "approval_3", "target": "end_1" }
  ]
}
```

---

## 2. 动态表单方案

### 方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **React Hook Form + JSON Schema** | 灵活、性能好 | 需要自己渲染 | ⭐⭐⭐⭐ |
| **react-jsonschema-form** | 开箱即用 | 样式定制复杂 | ⭐⭐⭐ |
| **Formily** | 阿里系、功能强 | 包体积大、学习曲线 | ⭐⭐⭐ |
| **自研** | 完全可控 | 周期长 | ⭐⭐ |

### 推荐方案：React Hook Form + 自研渲染器

```typescript
// 表单字段类型
type FieldType =
  | 'text'      // 单行文本
  | 'textarea'  // 多行文本
  | 'number'    // 数字
  | 'date'      // 日期
  | 'select'    // 下拉选择
  | 'radio'     // 单选
  | 'checkbox'  // 多选
  | 'file'      // 文件上传
  | 'table';    // 明细表（多行）

// 字段定义
interface FormField {
  name: string;
  type: FieldType;
  label: string;
  required: boolean;
  options?: { label: string; value: any }[];  // select/radio/checkbox用
  placeholder?: string;
  rules?: ValidationRule[];
  children?: FormField[];  // 明细表用
}

// 模板表单定义示例
{
  "fields": [
    { "name": "supplier_id", "type": "select", "label": "供应商", "required": true, "options": [] },
    { "name": "category_id", "type": "select", "label": "费用科目", "required": true },
    { "name": "items", "type": "table", "label": "采购明细", "children": [
      { "name": "product_name", "type": "text", "label": "产品名称" },
      { "name": "quantity", "type": "number", "label": "数量" },
      { "name": "unit_price", "type": "number", "label": "单价" },
      { "name": "amount", "type": "number", "label": "金额" }
    ]},
    { "name": "total_amount", "type": "number", "label": "总金额", "required": true },
    { "name": "attachments", "type": "file", "label": "附件", "multiple": true }
  ]
}
```

### 明细表数据结构

```typescript
// 表单提交数据
interface FormData {
  supplier_id: number;
  category_id: number;
  items: Array<{
    product_name: string;
    quantity: number;
    unit_price: number;
    amount: number;
  }>;
  total_amount: number;
  attachments: File[];
}
```

---

## 3. 文件存储方案

### 方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **MinIO** | S3兼容、自托管、简单 | 需要额外部署 | ⭐⭐⭐⭐ |
| **本地存储** | 无需额外服务 | 不适合分布式、备份麻烦 | ⭐⭐ |
| **阿里云 OSS** | 托管简单 | 费用、依赖第三方 | ⭐⭐⭐ |
| **七牛云** | 国内访问好 | 需实名认证 | ⭐⭐⭐ |

### 推荐方案：MinIO（开发/测试）+ S3兼容云存储（生产）

```yaml
# MinIO 配置
minio:
  endpoint: "localhost:9000"
  access_key: "minioadmin"
  secret_key: "minioadmin"
  bucket: "workflow-attachments"
  use_ssl: false
```

### 文件上传流程

```
1. 前端上传文件到 /api/v1/attachments/upload
2. 后端保存到 MinIO
3. 返回 file_path (UUID)
4. 存入 attachment 表
5. 预览/下载通过 /api/v1/attachments/:id/preview
```

### 支持的文件类型

| 类型 | MIME |
|------|------|
| PDF | application/pdf |
| 图片 | image/jpeg, image/png, image/gif |
| Excel | application/vnd.openxmlformats-officedocument.spreadsheetml.sheet |
| CSV | text/csv |

---

## 4. 条件表达式引擎

### 方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **govaluate** | 轻量、表达式丰富 | 不支持复杂逻辑 | ⭐⭐⭐⭐ |
| **道具** | 强大、AST | 学习曲线 | ⭐⭐⭐ |
| **expr** | 性能好、Go原生 | 较新 | ⭐⭐⭐ |
| **自定义解析** | 完全可控 | 开发周期长 | ⭐⭐ |

### 推荐方案：govaluate

```go
import "github.com/Knetic/govaluate"

// 使用示例
expression, _ := govaluate.NewEvaluableExpression("amount > 200000 && category == 'purchase'")

parameters := make(map[string]interface{})
parameters["amount"] = 250000
parameters["category"] = "purchase"

result, _ := expression.Evaluate(parameters)
// result = true
```

### 条件节点解析

```go
// 条件节点数据结构
type ConditionNode struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"` // >, <, ==, !=, >=, <=, in, contains
    Value    interface{} `json:"value"`
}

// 转换为 govaluate 表达式
func (c *ConditionNode) ToExpression() string {
    switch c.Operator {
    case ">":
        return fmt.Sprintf("%s > %v", c.Field, c.Value)
    case "<":
        return fmt.Sprintf("%s < %v", c.Field, c.Value)
    case "==":
        return fmt.Sprintf("%s == '%v'", c.Field, c.Value)
    case ">=":
        return fmt.Sprintf("%s >= %v", c.Field, c.Value)
    case "<=":
        return fmt.Sprintf("%s <= %v", c.Field, c.Value)
    case "!=":
        return fmt.Sprintf("%s != '%v'", c.Field, c.Value)
    case "in":
        return fmt.Sprintf("['%v'].contains(%s)", c.Value, c.Field)
    case "contains":
        return fmt.Sprintf("%s.contains('%v')", c.Field, c.Value)
    default:
        return "true"
    }
}
```

---

## 5. 技术栈汇总

### 后端

| 技术 | 选型 | 理由 |
|------|------|------|
| **语言** | Go 1.26+ | 性能好、并发支持 |
| **框架** | Gin | 轻量、社区活跃 |
| **ORM** | GORM | Go标配、功能全 |
| **数据库** | PostgreSQL | JSONB支持好、复杂查询强 |
| **缓存** | Redis | 会话、队列 |
| **文件存储** | MinIO (S3兼容) | 自托管、简单 |
| **表达式** | govaluate | 轻量、条件判断 |
| **认证** | JWT | 成熟、无状态 |

### 前端

| 技术 | 选型 | 理由 |
|------|------|------|
| **框架** | React 18 + TypeScript | 类型安全、生态好 |
| **构建** | Vite | 快速、HMR好 |
| **UI库** | Ant Design 5.x | 组件丰富 |
| **状态** | Redux Toolkit | 多人协作可控 |
| **流程设计** | React Flow | 开源、定制灵活 |
| **表单** | React Hook Form | 性能好、灵活 |
| **HTTP** | Axios | 拦截器强 |
| **路由** | React Router 6 | 路由标配 |

### 部署

| 技术 | 用途 |
|------|------|
| **Docker** | 容器化 |
| **Docker Compose** | 本地开发 |
| **Nginx** | 反向代理 |

---

## 6. 项目结构

```
C:/project/Workflow_claude/
├── backend/
│   ├── cmd/server/
│   │   └── main.go
│   ├── config/
│   ├── internal/
│   │   ├── domain/          # 领域实体
│   │   │   ├── company/
│   │   │   ├── department/
│   │   │   ├── employee/
│   │   │   ├── supplier/
│   │   │   ├── expense_category/
│   │   │   ├── workflow/
│   │   │   ├── instance/
│   │   │   ├── task/
│   │   │   ├── attachment/
│   │   │   └── notification/
│   │   ├── repository/       # 数据访问
│   │   ├── service/        # 业务逻辑
│   │   │   ├── engine/     # 流程引擎
│   │   │   ├── approval/
│   │   │   └── auth/
│   │   ├── handler/        # HTTP处理
│   │   │   └── api/v1/
│   │   └── pkg/
│   │       ├── jwt/
│   │       ├── email/
│   │       ├── storage/    # MinIO
│   │       └── expression/ # govaluate
│   ├── migrations/
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   │   ├── workflow/  # 流程设计器
│   │   │   ├── form/       # 表单渲染器
│   │   │   └── common/
│   │   ├── pages/
│   │   ├── store/
│   │   ├── types/
│   │   ├── hooks/
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
├── docs/
│   └── design/
├── docker-compose.yml
└── README.md
```

---

## 7. 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| React Flow 定制复杂 | 中 | 中 | 提前研究文档，从简单节点开始 |
| 动态表单验证复杂 | 中 | 中 | 分阶段开发，先做基础字段 |
| MinIO 运维成本 | 低 | 低 | 开发用Docker，生产可选云存储 |
| 条件表达式边界情况 | 中 | 高 | 充分单元测试，覆盖边界值 |

---

## 8. 下一步

1. ✅ 技术选型确定
2. ⬜ 创建数据库迁移脚本
3. ⬜ 实现后端 CRUD API
4. ⬜ 实现流程引擎
5. ⬜ 前端页面开发
6. ⬜ 集成测试

---

**结论：技术选型可行，推荐采用上述方案。**

---

## 9. 可移植性设计

### 部署目标支持

| 部署目标 | 支持情况 | 说明 |
|----------|----------|------|
| 本地机器 | ✅ | Docker Compose 一键启动 |
| 云端 VM | ✅ | 同样的 Docker 部署 |
| 机房实体服务器 | ✅ | Linux 服务器 + Docker |
| 云托管服务 | ✅ | RDS + ElastiCache + S3 |

### 统一存储接口

后端使用 S3 兼容 API，通过环境变量切换存储 provider：

```
┌─────────────────────────────────────┐
│          统一存储接口                 │
│  (后端代码不变，只换配置)             │
└─────────────────────────────────────┘
         ↓            ↓            ↓
    ┌────────┐  ┌─────────┐  ┌─────────┐
    │ MinIO  │  │ AWS S3  │  │阿里云OSS│
    │(自托管)│  │(云托管) │  │(云托管) │
    └────────┘  └─────────┘  └─────────┘
```

### 环境变量配置

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=workflow
DB_PASSWORD=xxx

# Redis
REDIS_HOST=localhost
REDIS_PASSWORD=xxx

# 存储 (可选: minio, s3, oss)
STORAGE_TYPE=minio
MINIO_ENDPOINT=localhost:9000
MINIO_BUCKET=workflow

# 或云端 S3
# STORAGE_TYPE=s3
# AWS_REGION=ap-east-1
# AWS_S3_BUCKET=workflow-prod

# 或阿里云 OSS
# STORAGE_TYPE=oss
# OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com

# JWT
JWT_SECRET=xxx

# LDAP (后续)
LDAP_HOST=ldap://xxx
```

### Docker Compose 架构

```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - STORAGE_TYPE=${STORAGE_TYPE:-minio}
      - MINIO_ENDPOINT=minio:9000
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
      - minio
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=workflow
      - POSTGRES_USER=workflow
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    command: redis-server --requirepass ${REDIS_PASSWORD}
    restart: unless-stopped

  minio:
    image: minio/minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=${MINIO_USER:-minioadmin}
      - MINIO_ROOT_PASSWORD=${MINIO_PASSWORD:-minioadmin}
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - frontend
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  minio_data:
```

### Nginx 配置

```nginx
server {
    listen 80;
    server_name _;

    # 前端静态资源
    location / {
        proxy_pass http://frontend:3000;
    }

    # API 代理
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # MinIO 控制台 (开发用)
    location /minio/ {
        proxy_pass http://minio:9001/;
    }
}
```

### 迁移能力

| 场景 | 方案 |
|------|------|
| 本地 → 云端 VM | 直接迁移 Docker Compose，配置文件一起搬 |
| VM → 机房实体机 | 同上，Linux 环境一致 |
| 自托管 → 云托管 | 改环境变量指向 RDS/ElastiCache/S3，无需代码改动 |
| 数据迁移 | PostgreSQL 和 S3 数据都有导出工具 |

### 快速启动脚本

```bash
#!/bin/bash
# 一键启动所有服务

# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
vim .env

# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

**可移植性设计完成，支持任意环境部署。**
