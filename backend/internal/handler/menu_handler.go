package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuService: service.NewMenuService(),
	}
}

// List 菜单树
// @Summary 菜单树
// @Tags 菜单管理
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysMenu}
// @Router /api/system/menu [get]
func (h *MenuHandler) List(c *gin.Context) {
	tree, err := h.menuService.SelectTree()
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}
	pkg.OK(c, tree)
}

// SelectAll 查询所有菜单(平铺)
// @Summary 查询所有菜单
// @Tags 菜单管理
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysMenu}
// @Router /api/system/menu/all [get]
func (h *MenuHandler) SelectAll(c *gin.Context) {
	menus, err := h.menuService.SelectAll()
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}
	pkg.OK(c, menus)
}

// Get 获取菜单详情
// @Summary 获取菜单详情
// @Tags 菜单管理
// @Security Bearer
// @Param id path int true "菜单ID"
// @Success 200 {object} pkg.Response{data=model.SysMenu}
// @Router /api/system/menu/{id} [get]
func (h *MenuHandler) Get(c *gin.Context) {
	menuID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	menu, err := h.menuService.GetMenuByID(menuID)
	if err != nil {
		pkg.Fail(c, "菜单不存在")
		return
	}

	pkg.OK(c, menu)
}

// Create 创建菜单
// @Summary 创建菜单
// @Tags 菜单管理
// @Security Bearer
// @Accept json
// @Param body body model.CreateMenuRequest true "菜单信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/menu [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req model.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.menuService.Create(&req); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "创建成功")
}

// Update 更新菜单
// @Summary 更新菜单
// @Tags 菜单管理
// @Security Bearer
// @Accept json
// @Param id path int true "菜单ID"
// @Param body body model.CreateMenuRequest true "菜单信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/menu/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	menuID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	var req model.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.menuService.Update(menuID, &req); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "更新成功")
}

// Delete 删除菜单
// @Summary 删除菜单
// @Tags 菜单管理
// @Security Bearer
// @Param id path int true "菜单ID"
// @Success 200 {object} pkg.Response
// @Router /api/system/menu/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	menuID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	if err := h.menuService.Delete(menuID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "删除成功")
}
