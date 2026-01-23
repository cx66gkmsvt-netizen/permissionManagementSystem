package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// List 角色列表
// @Summary 角色列表
// @Tags 角色管理
// @Security Bearer
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} pkg.Response{data=model.PageResult}
// @Router /api/system/role [get]
func (h *RoleHandler) List(c *gin.Context) {
	var query model.RoleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	result, err := h.roleService.List(&query)
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}

	pkg.OK(c, result)
}

// SelectAll 查询所有角色
// @Summary 查询所有角色
// @Tags 角色管理
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysRole}
// @Router /api/system/role/all [get]
func (h *RoleHandler) SelectAll(c *gin.Context) {
	roles, err := h.roleService.SelectAll()
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}
	pkg.OK(c, roles)
}

// Get 获取角色详情
// @Summary 获取角色详情
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} pkg.Response{data=model.SysRole}
// @Router /api/system/role/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	role, err := h.roleService.GetRoleByID(roleID)
	if err != nil {
		pkg.Fail(c, "角色不存在")
		return
	}

	// 获取菜单ID和部门ID
	menuIDs, _ := h.roleService.GetRoleMenuIDs(roleID)
	deptIDs, _ := h.roleService.GetRoleDeptIDs(roleID)

	pkg.OK(c, gin.H{
		"role":    role,
		"menuIds": menuIDs,
		"deptIds": deptIDs,
	})
}

// Create 创建角色
// @Summary 创建角色
// @Tags 角色管理
// @Security Bearer
// @Accept json
// @Param body body model.CreateRoleRequest true "角色信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/role [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.roleService.Create(c.Request.Context(), &req); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "创建成功")
}

// Update 更新角色
// @Summary 更新角色
// @Tags 角色管理
// @Security Bearer
// @Accept json
// @Param id path int true "角色ID"
// @Param body body model.UpdateRoleRequest true "角色信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/role/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.roleService.Update(c.Request.Context(), roleID, &req); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "更新成功")
}

// Delete 删除角色
// @Summary 删除角色
// @Tags 角色管理
// @Security Bearer
// @Param id path int true "角色ID"
// @Success 200 {object} pkg.Response
// @Router /api/system/role/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.roleService.Delete(c.Request.Context(), roleID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "删除成功")
}
