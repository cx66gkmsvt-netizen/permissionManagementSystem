package handler

import (
	"github.com/gin-gonic/gin"

	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/repository"
)

type ProfileHandler struct {
	userRepo *repository.UserRepository
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		userRepo: repository.NewUserRepository(),
	}
}

// GetProfile 获取个人信息
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		pkg.Fail(c, "获取个人信息失败")
		return
	}
	pkg.OK(c, user)
}

// UpdateProfile 更新个人信息
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	db := repository.GetDB()
	updates := map[string]interface{}{
		"nick_name": req.NickName,
		"phone":     req.Phone,
		"email":     req.Email,
	}

	if err := db.Table("sys_user").Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		pkg.Fail(c, "更新失败")
		return
	}

	pkg.OKMsg(c, "更新成功")
}

// UpdatePassword 修改密码
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req model.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	// 获取当前用户
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		pkg.Fail(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !pkg.CheckPassword(req.OldPassword, user.Password) {
		pkg.Fail(c, "旧密码错误")
		return
	}

	// 更新新密码
	newPwd, _ := pkg.HashPassword(req.NewPassword)
	db := repository.GetDB()
	if err := db.Table("sys_user").Where("user_id = ?", userID).Update("password", newPwd).Error; err != nil {
		pkg.Fail(c, "修改密码失败")
		return
	}

	pkg.OKMsg(c, "密码修改成功")
}
