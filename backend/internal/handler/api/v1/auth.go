package v1

import (
	"net/http"

	"workflow-system/internal/pkg/jwt"
	"workflow-system/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	jwt           *jwt.JWT
	employeeSvc   *service.EmployeeService
}

func NewAuthHandler(jwt *jwt.JWT, employeeSvc *service.EmployeeService) *AuthHandler {
	return &AuthHandler{jwt: jwt, employeeSvc: employeeSvc}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock 登录（用于测试）
	if req.Username == "admin" && req.Password == "admin123" {
		token, _ := h.jwt.GenerateToken(2, "admin", 1)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":         2,
				"username":   "admin",
				"company_id": 1,
			},
		})
		return
	}

	// 查询用户进行真实验证
	emp, err := h.employeeSvc.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// 验证密码
	if !h.employeeSvc.VerifyPassword(emp, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password incorrect"})
		return
	}

	// 生成 Token
	token, err := h.jwt.GenerateToken(emp.ID, emp.Username, emp.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         emp.ID,
			"username":   emp.Username,
			"company_id": emp.CompanyID,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Current(c *gin.Context) {
	if userID, exists := c.Get("user_id"); exists {
		c.JSON(http.StatusOK, gin.H{
			"id":       userID,
			"username": c.GetString("username"),
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
}
