-- Workflow 系统初始数据库 Schema

-- 公司表
CREATE TABLE IF NOT EXISTS company (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    short_name VARCHAR(100),
    status SMALLINT DEFAULT 1,
    schema_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 部门表
CREATE TABLE IF NOT EXISTS department (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES department(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    leader_id BIGINT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 部门审批链
CREATE TABLE IF NOT EXISTS department_approval_chain (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES department(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL,
    step_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 员工表
CREATE TABLE IF NOT EXISTS employee (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    level VARCHAR(50),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 员工部门关联
CREATE TABLE IF NOT EXISTS employee_department (
    employee_id BIGINT NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    department_id BIGINT NOT NULL REFERENCES department(id) ON DELETE CASCADE,
    PRIMARY KEY (employee_id, department_id)
);

-- 员工银行账户
CREATE TABLE IF NOT EXISTS employee_bank_account (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    bank_name VARCHAR(200) NOT NULL,
    bank_branch VARCHAR(200),
    bank_account VARCHAR(100) NOT NULL,
    account_holder VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 供应商表
CREATE TABLE IF NOT EXISTS supplier (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    contact VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    address VARCHAR(500),
    bank_name VARCHAR(200),
    bank_account VARCHAR(100),
    tax_number VARCHAR(50),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 费用科目表
CREATE TABLE IF NOT EXISTS expense_category (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    parent_id BIGINT REFERENCES expense_category(id) ON DELETE SET NULL,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 流程定义表
CREATE TABLE IF NOT EXISTS workflow_definition (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    version INT DEFAULT 1,
    graph_data JSONB DEFAULT '{}',
    form_fields JSONB DEFAULT '[]',
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 流程实例表
CREATE TABLE IF NOT EXISTS workflow_instance (
    id BIGSERIAL PRIMARY KEY,
    definition_id BIGINT NOT NULL REFERENCES workflow_definition(id),
    company_id BIGINT NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    initiator_id BIGINT NOT NULL REFERENCES employee(id),
    form_data JSONB DEFAULT '{}',
    status SMALLINT DEFAULT 1,
    current_nodes JSONB DEFAULT '[]',
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 审批任务表
CREATE TABLE IF NOT EXISTS approval_task (
    id BIGSERIAL PRIMARY KEY,
    instance_id BIGINT NOT NULL REFERENCES workflow_instance(id) ON DELETE CASCADE,
    node_id VARCHAR(50) NOT NULL,
    node_name VARCHAR(100) NOT NULL,
    assignee_id BIGINT NOT NULL REFERENCES employee(id),
    status SMALLINT DEFAULT 1,
    action VARCHAR(20),
    comment TEXT,
    completed_at TIMESTAMP,
    version INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 附件表
CREATE TABLE IF NOT EXISTS attachment (
    id BIGSERIAL PRIMARY KEY,
    instance_id BIGINT NOT NULL REFERENCES workflow_instance(id) ON DELETE CASCADE,
    field_name VARCHAR(100),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 通知表
CREATE TABLE IF NOT EXISTS notification (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    type VARCHAR(50),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT,
    user_id BIGINT,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(100),
    details JSONB,
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_department_company ON department(company_id);
CREATE INDEX IF NOT EXISTS idx_employee_company ON employee(company_id);
CREATE INDEX IF NOT EXISTS idx_supplier_company ON supplier(company_id);
CREATE INDEX IF NOT EXISTS idx_expense_category_company ON expense_category(company_id);
CREATE INDEX IF NOT EXISTS idx_workflow_definition_company ON workflow_definition(company_id);
CREATE INDEX IF NOT EXISTS idx_workflow_instance_company ON workflow_instance(company_id);
CREATE INDEX IF NOT EXISTS idx_workflow_instance_definition ON workflow_instance(definition_id);
CREATE INDEX IF NOT EXISTS idx_approval_task_instance ON approval_task(instance_id);
CREATE INDEX IF NOT EXISTS idx_approval_task_assignee ON approval_task(assignee_id);
CREATE INDEX IF NOT EXISTS idx_notification_user ON notification(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_company ON audit_log(company_id);

-- 插入默认公司
INSERT INTO company (code, name, short_name, status, schema_name)
VALUES ('DEMO', '演示公司', 'DEMO', 1, 'demo')
ON CONFLICT (code) DO NOTHING;
