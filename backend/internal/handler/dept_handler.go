package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"user-center/internal/middleware"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"
)

type DeptHandler struct {
	deptService *service.DeptService
}

func NewDeptHandler() *DeptHandler {
	return &DeptHandler{
		deptService: service.NewDeptService(),
	}
}

// List 部门树
// @Summary 部门树
// @Tags 部门管理
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysDept}
// @Router /api/system/dept [get]
func (h *DeptHandler) List(c *gin.Context) {
	tree, err := h.deptService.SelectTree()
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}
	pkg.OK(c, tree)
}

// SelectAll 查询所有部门(平铺)
// @Summary 查询所有部门
// @Tags 部门管理
// @Security Bearer
// @Success 200 {object} pkg.Response{data=[]model.SysDept}
// @Router /api/system/dept/all [get]
func (h *DeptHandler) SelectAll(c *gin.Context) {
	depts, err := h.deptService.SelectAll()
	if err != nil {
		pkg.Fail(c, "查询失败")
		return
	}
	pkg.OK(c, depts)
}

// Get 获取部门详情
// @Summary 获取部门详情
// @Tags 部门管理
// @Security Bearer
// @Param id path int true "部门ID"
// @Success 200 {object} pkg.Response{data=model.SysDept}
// @Router /api/system/dept/{id} [get]
func (h *DeptHandler) Get(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	dept, err := h.deptService.GetDeptByID(deptID)
	if err != nil {
		pkg.Fail(c, "部门不存在")
		return
	}

	pkg.OK(c, dept)
}

// Create 创建部门
// @Summary 创建部门
// @Tags 部门管理
// @Security Bearer
// @Accept json
// @Param body body model.CreateDeptRequest true "部门信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/dept [post]
func (h *DeptHandler) Create(c *gin.Context) {
	var req model.CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	// 90
	operatorID := middleware.GetUserID(c)
	if err := h.deptService.Create(&req, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "创建成功")
}

// Update 更新部门
// @Summary 更新部门
// @Tags 部门管理
// @Security Bearer
// @Accept json
// @Param id path int true "部门ID"
// @Param body body model.CreateDeptRequest true "部门信息"
// @Success 200 {object} pkg.Response
// @Router /api/system/dept/{id} [put]
func (h *DeptHandler) Update(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	var req model.CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := h.deptService.Update(deptID, &req, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "更新成功")
}

// Delete 删除部门
// @Summary 删除部门
// @Tags 部门管理
// @Security Bearer
// @Param id path int true "部门ID"
// @Success 200 {object} pkg.Response
// @Router /api/system/dept/{id} [delete]
func (h *DeptHandler) Delete(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		pkg.FailCode(c, pkg.INVALID_PARAM, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	if err := h.deptService.Delete(deptID, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	pkg.OKMsg(c, "删除成功")
}
