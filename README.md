# Workflow 审批系统

企业级多公司审批流程管理系统。

## 技术栈

- **后端**: Go 1.26 + Gin + GORM + PostgreSQL
- **前端**: React 18 + TypeScript + Vite + Ant Design
- **流程引擎**: React Flow
- **部署**: Docker + Docker Compose + Nginx

## 快速开始

### 方式一: Docker (推荐)

```bash
# 启动所有服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
```

访问 http://localhost

### 方式二: 本地开发

```bash
# 后端
cd backend
cp .env.example .env
go mod download
go run cmd/server/main.go

# 前端 (新终端)
cd frontend
npm install
npm run dev
```

## 默认账号

- 用户名: admin
- 密码: admin123
- 公司代码: DEMO

## API 文档

启动服务后访问 http://localhost/api/v1 下的各端点。

## 开发指南

详见 [docs/](docs/)
