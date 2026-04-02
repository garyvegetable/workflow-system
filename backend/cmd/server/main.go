package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"workflow-system/internal/handler/api/v1"
	"workflow-system/internal/pkg/database"
	"workflow-system/internal/pkg/email"
	"workflow-system/internal/pkg/jwt"
	"workflow-system/internal/repository"
	"workflow-system/internal/service"
	"workflow-system/internal/service/engine"
	"workflow-system/internal/service/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	godotenv.Load()

	// 数据库连接
	db := database.InitDB()

	// 初始化仓库
	companyRepo := repository.NewCompanyRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	supRepo := repository.NewSupplierRepository(db)
	expenseRepo := repository.NewExpenseCategoryRepository(db)
	bankRepo := repository.NewEmployeeBankAccountRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	instanceRepo := repository.NewInstanceRepository(db)
	positionRepo := repository.NewPositionRepository(db)
	systemSettingRepo := repository.NewSystemSettingRepository(db)

	// 初始化服务
	companyService := service.NewCompanyService(companyRepo)
	deptService := service.NewDepartmentService(deptRepo, empRepo)
	empService := service.NewEmployeeService(empRepo, deptRepo, bankRepo)
	supService := service.NewSupplierService(supRepo)
	expenseService := service.NewExpenseCategoryService(expenseRepo)
	workflowService := service.NewWorkflowService(workflowRepo, instanceRepo)
	positionService := service.NewPositionService(positionRepo)
	systemSettingService := service.NewSystemSettingService(systemSettingRepo)

	// 邮件服务（需要系统设置）
	emailSvc := email.NewEmailService()
	// 尝试加载 SMTP 配置
	if smtpConfig, err := systemSettingService.GetSMTPConfig(); err == nil {
		port := 587
		if p, ok := smtpConfig["smtp_port"]; ok {
			if parsed, err := strconv.Atoi(p); err == nil {
				port = parsed
			}
		}
		emailSvc.Configure(
			smtpConfig["smtp_host"],
			port,
			smtpConfig["smtp_user"],
			smtpConfig["smtp_password"],
			smtpConfig["smtp_from"],
		)
		log.Println("SMTP 配置已加载")
	} else {
		log.Println("SMTP 未配置，将使用 Mock 模式")
	}

	// 通知
	notifRepo := repository.NewNotificationRepository(db)
	notifService := service.NewNotificationService(notifRepo, emailSvc)

	// 初始化流程引擎
	wfEngine := engine.NewWorkflowEngine(engine.WorkflowEngineDeps{
		WorkflowRepo:   workflowRepo,
		InstanceRepo:   instanceRepo,
		DepartmentRepo: deptRepo,
		NotifService:   notifService,
	})
	approvalService := service.NewApprovalService(instanceRepo, workflowRepo, wfEngine)
	approvalHandler := v1.NewApprovalHandler(approvalService, empRepo)

	// 启动超时调度器（每5分钟检查一次）
	timeoutScheduler := scheduler.NewTimeoutScheduler(instanceRepo, wfEngine, 5*time.Minute)
	timeoutScheduler.Start()

	// JWT
	jwt.Init(os.Getenv("JWT_SECRET"))

	// 路由
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API v1
	v1Router := r.Group("/api/v1")
	{
		// 认证（公开）
		authHandler := v1.NewAuthHandler(jwt.NewJWT(), empService)
		v1Router.POST("/auth/login", authHandler.Login)
		v1Router.POST("/auth/logout", authHandler.Logout)

		// 需要认证的路由
		auth := v1Router.Group("")
		auth.Use(v1.JWTAuthMiddleware())
		{
			auth.GET("/auth/current", authHandler.Current)

			// 公司
			companyHandler := v1.NewCompanyHandler(companyService)
			auth.GET("/companies", companyHandler.List)
			auth.POST("/companies", companyHandler.Create)
			auth.GET("/companies/:id", companyHandler.Get)
			auth.PUT("/companies/:id", companyHandler.Update)
			auth.DELETE("/companies/:id", companyHandler.Delete)

			// 职位
			positionHandler := v1.NewPositionHandler(positionService)
			auth.GET("/positions", positionHandler.List)
			auth.POST("/positions", positionHandler.Create)
			auth.GET("/positions/:id", positionHandler.Get)
			auth.PUT("/positions/:id", positionHandler.Update)
			auth.DELETE("/positions/:id", positionHandler.Delete)

			// 系统设置
			systemSettingHandler := v1.NewSystemSettingHandler(systemSettingService)
			auth.GET("/system-settings", systemSettingHandler.List)
			auth.GET("/system-settings/:key", systemSettingHandler.Get)
			auth.PUT("/system-settings/:key", systemSettingHandler.Set)
			auth.DELETE("/system-settings/:key", systemSettingHandler.Delete)

			// 部门
			deptHandler := v1.NewDepartmentHandler(deptService)
			auth.GET("/departments", deptHandler.List)
			auth.POST("/departments", deptHandler.Create)
			auth.GET("/departments/:id", deptHandler.Get)
			auth.PUT("/departments/:id", deptHandler.Update)
			auth.DELETE("/departments/:id", deptHandler.Delete)
			auth.GET("/departments/:id/approval-chain", deptHandler.GetApprovalChain)
			auth.PUT("/departments/:id/approval-chain", deptHandler.SetApprovalChain)

			// 员工
			empHandler := v1.NewEmployeeHandler(empService, emailSvc)
			auth.GET("/employees", empHandler.List)
			auth.GET("/employees/search", empHandler.Search)
			auth.POST("/employees", empHandler.Create)
			auth.GET("/employees/:id", empHandler.Get)
			auth.PUT("/employees/:id", empHandler.Update)
			auth.DELETE("/employees/:id", empHandler.Delete)
			auth.GET("/employees/:id/bank-accounts", empHandler.ListBankAccounts)
			auth.POST("/employees/:id/bank-accounts", empHandler.CreateBankAccount)
			auth.PUT("/employees/:id/bank-accounts/:aid", empHandler.UpdateBankAccount)
			auth.DELETE("/employees/:id/bank-accounts/:aid", empHandler.DeleteBankAccount)

			// 供应商
			supHandler := v1.NewSupplierHandler(supService)
			auth.GET("/suppliers", supHandler.List)
			auth.POST("/suppliers", supHandler.Create)
			auth.GET("/suppliers/:id", supHandler.Get)
			auth.PUT("/suppliers/:id", supHandler.Update)
			auth.DELETE("/suppliers/:id", supHandler.Delete)

			// 费用科目
			expHandler := v1.NewExpenseCategoryHandler(expenseService)
			auth.GET("/expense-categories", expHandler.List)
			auth.POST("/expense-categories", expHandler.Create)
			auth.GET("/expense-categories/:id", expHandler.Get)
			auth.PUT("/expense-categories/:id", expHandler.Update)
			auth.DELETE("/expense-categories/:id", expHandler.Delete)

			// 流程
			wfHandler := v1.NewWorkflowHandler(workflowService)
			auth.GET("/workflows", wfHandler.List)
			auth.POST("/workflows", wfHandler.Create)
			auth.GET("/workflows/:id", wfHandler.Get)
			auth.PUT("/workflows/:id", wfHandler.Update)
			auth.POST("/workflows/:id/publish", wfHandler.Publish)
			auth.POST("/workflows/:id/disable", wfHandler.Disable)
			auth.POST("/workflows/:id/copy", wfHandler.Copy)
			auth.DELETE("/workflows/:id", wfHandler.Delete)
			auth.POST("/workflows/instances", wfHandler.CreateInstance)
			auth.GET("/workflows/instances/my", wfHandler.GetMyApplications)
			auth.GET("/workflows/instances/:id", wfHandler.GetInstance)
			auth.POST("/workflows/instances/:id/cancel", wfHandler.CancelInstance)

			// 审批
			auth.GET("/tasks/pending", approvalHandler.ListPending)
			auth.GET("/tasks/handled", approvalHandler.ListHandled)
			auth.POST("/tasks/:id/approve", approvalHandler.Approve)
			auth.POST("/tasks/:id/reject", approvalHandler.Reject)
			auth.POST("/tasks/:id/transfer", approvalHandler.Transfer)
			auth.GET("/tasks/:id/history", approvalHandler.GetHistory)
			auth.POST("/tasks/batch-approve", approvalHandler.BatchApprove)
			auth.POST("/tasks/batch-reject", approvalHandler.BatchReject)
			auth.POST("/tasks/:id/add-approver", approvalHandler.AddApprover)
			auth.POST("/tasks/:id/remove-approver", approvalHandler.RemoveApprover)

			// 附件
			attachHandler := v1.NewAttachmentHandler()
			auth.POST("/attachments/upload", attachHandler.Upload)
			auth.GET("/attachments/:id/download", attachHandler.Download)
			auth.GET("/attachments/:id/preview", attachHandler.Preview)

			// 通知
			notifHandler := v1.NewNotificationHandler(notifService)
			auth.GET("/notifications", notifHandler.List)
			auth.PUT("/notifications/:id/read", notifHandler.MarkRead)
			auth.PUT("/notifications/read-all", notifHandler.MarkAllRead)

			// 审计日志
			auditHandler := v1.NewAuditHandler()
			auth.GET("/audit-logs", auditHandler.List)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(r.Run(":" + port))
}
