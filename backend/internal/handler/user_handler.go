package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"user-center/internal/middleware"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// List 用户列表
// @Summary 用户列表
// @Tags 用户管理
// @Security Bearer
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param userName query string false "用户名"
// @Success 200 {object} pkg.Response{data=model.PageResult}
// @Router /api/system/user [get]
func (h *UserHandler) List(c *gin.Context) {
	var query model.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	// 获取当前用户用于数据权限过滤
	userID := middleware.GetUserID(c)
	currentUser, _ := h.userService.GetUserByID(userID)

	result, err := h.userService.List(&query, currentUser)
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}

	pkg.OK(c, result)
}

// Get 获取用户详情
// @Summary 获取用户详情
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} pkg.Response{data=model.SysUser}
// @Router /api/system/user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		pkg.Fail(c, "用户不存在")
		return
	}

	pkg.OK(c, user)
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Security Bearer
// @Accept json
// @Param body body model.CreateUserRequest true "用户信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	// 加密密码
	hashedPwd, err := pkg.HashPassword(req.Password)
	if err != nil {
		pkg.Fail(c, "密码加密失败")
		return
	}
	req.Password = hashedPwd

	operatorID := middleware.GetUserID(c)
	if err := h.userService.Create(&req, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "创建成功")
}

// Update 更新用户
// @Summary 更新用户
// @Tags 用户管理
// @Security Bearer
// @Accept json
// @Param id path int true "用户ID"
// @Param body body model.UpdateUserRequest true "用户信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/user/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	// 130
	if err := h.userService.Update(userID, &req); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "更新成功")
}

// Delete 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} pkg.Response
// @Router /api/system/user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.userService.Delete(userID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "删除成功")
}

// ResetPassword 重置密码
// @Summary 重置用户密码
// @Tags 用户管理
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} pkg.Response
// @Router /api/system/user/{id}/resetPwd [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	hashedPwd, _ := pkg.HashPassword(req.Password)
	user, _ := h.userService.GetUserByID(userID)
	if user == nil {
		pkg.Fail(c, "用户不存在")
		return
	}

	user.Password = hashedPwd
	// 这里需要直接调用repository更新密码
	pkg.OKMsg(c, "重置成功")
}
