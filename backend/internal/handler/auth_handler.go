package handler

import (
	"github.com/gin-gonic/gin"

	"user-center/internal/config"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(cfg),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "登录信息"
// @Success 200 {object} pkg.Response{data=model.LoginResponse}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	resp, err := h.authService.Login(&req, c.ClientIP())
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OK(c, resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Tags 认证
// @Success 200 {object} pkg.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Token 可以在 Redis 中加入黑名单
	pkg.OKMsg(c, "登出成功")
}

// GetUserInfo 获取用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Security Bearer
// @Success 200 {object} pkg.Response{data=model.UserInfo}
// @Router /api/auth/info [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt64("userID")
	info, err := h.authService.GetUserInfo(userID)
	if err != nil {
		pkg.Fail(c, "获取用户信息失败")
		return
	}
	pkg.OK(c, info)
}

// GetRoutes 获取路由菜单
// @Summary 获取当前用户路由
// @Tags 认证
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysMenu}
// @Router /api/auth/routes [get]
func (h *AuthHandler) GetRoutes(c *gin.Context) {
	userID := c.GetInt64("userID")
	routes, err := h.authService.GetRoutes(userID)
	if err != nil {
		pkg.Fail(c, "获取路由失败")
		return
	}
	pkg.OK(c, routes)
}
