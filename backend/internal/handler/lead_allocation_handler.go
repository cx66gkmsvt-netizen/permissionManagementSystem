package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type LeadAllocationHandler struct {
	svc *service.LeadAllocationService
}

func NewLeadAllocationHandler() *LeadAllocationHandler {
	return &LeadAllocationHandler{
		svc: service.NewLeadAllocationService(),
	}
}

// List 获取分配列表
func (h *LeadAllocationHandler) List(c *gin.Context) {
	var query model.LeadAllocationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	date := c.Query("date")
	if date == "" {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	result, err := h.svc.GetAllocationList(&query, date)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, result)
}

// Update 更新分配记录
func (h *LeadAllocationHandler) Update(c *gin.Context) {
	ccID, _ := strconv.ParseInt(c.Param("ccId"), 10, 64)
	date := c.Query("date")
	if date == "" {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	var dto model.LeadAllocationUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdateAllocation(ccID, date, &dto); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// BatchUpdateIsAllocated 批量更新是否分配
func (h *LeadAllocationHandler) BatchUpdateIsAllocated(c *gin.Context) {
	var dto model.LeadAllocationBatchUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	if err := h.svc.BatchUpdateIsAllocated(dto.CCIDs, dto.Date, dto.IsAllocated, dto.Reason); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Stats 获取分配统计
func (h *LeadAllocationHandler) Stats(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	stats, err := h.svc.GetAllocationStats(date)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, stats)
}

// Detail 获取CC分配详情历史
func (h *LeadAllocationHandler) Detail(c *gin.Context) {
	ccID, _ := strconv.ParseInt(c.Param("ccId"), 10, 64)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		pkg.FailCode(c, 400, "请提供日期范围")
		return
	}

	list, err := h.svc.GetAllocationDetail(ccID, startDate, endDate)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// SingleDetail 获取分配详情单日明细
func (h *LeadAllocationHandler) SingleDetail(c *gin.Context) {
	ccID, _ := strconv.ParseInt(c.Param("ccId"), 10, 64)
	date := c.Query("date")

	if date == "" {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	result, err := h.svc.GetSingleDetail(ccID, date)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, result)
}
