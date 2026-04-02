package domain

import (
	"workflow-system/internal/domain/attachment"
	"workflow-system/internal/domain/company"
	"workflow-system/internal/domain/department"
	"workflow-system/internal/domain/employee"
	"workflow-system/internal/domain/expense_category"
	"workflow-system/internal/domain/instance"
	"workflow-system/internal/domain/notification"
	"workflow-system/internal/domain/supplier"
	"workflow-system/internal/domain/task"
	"workflow-system/internal/domain/workflow"
)

// Re-export types from sub-packages for backwards compatibility
type Employee = employee.Employee
type EmployeeBankAccount = employee.EmployeeBankAccount
type EmployeeDepartment = employee.EmployeeDepartment
type Company = company.Company
type Department = department.Department
type DepartmentApprovalChain = department.DepartmentApprovalChain
type TreeNode = department.TreeNode
type Supplier = supplier.Supplier
type ExpenseCategory = expense_category.ExpenseCategory
type ExpenseCategoryTreeNode = expense_category.ExpenseCategoryTreeNode
type WorkflowDefinition = workflow.WorkflowDefinition
type WorkflowInstance = instance.WorkflowInstance
type ApprovalTask = task.ApprovalTask
type Attachment = attachment.Attachment
type Notification = notification.Notification
