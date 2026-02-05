package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{
		svc: service.NewAttendanceService(),
	}
}

// List 获取在班列表
func (h *AttendanceHandler) List(c *gin.Context) {
	var query model.CCQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	// 获取日期列表
	dates := c.QueryArray("dates")
	if len(dates) == 0 {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	list, err := h.svc.GetAttendanceList(&query, dates)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// Update 更新在班状态
func (h *AttendanceHandler) Update(c *gin.Context) {
	ccID, _ := strconv.ParseInt(c.Param("ccId"), 10, 64)
	date := c.Param("date")

	var dto model.AttendanceUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, _ := getOperatorInfo(c)
	if err := h.svc.UpdateAttendance(ccID, date, dto.Status, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// BatchUpdate 批量更新在班状态
func (h *AttendanceHandler) BatchUpdate(c *gin.Context) {
	var dto model.AttendanceBatchUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, _ := getOperatorInfo(c)
	if err := h.svc.BatchUpdateAttendance(dto.Items, operatorID); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Stats 获取在班统计
func (h *AttendanceHandler) Stats(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	stats, err := h.svc.GetAttendanceStats(date)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, stats)
}

// History 获取CC在班历史
func (h *AttendanceHandler) History(c *gin.Context) {
	ccID, _ := strconv.ParseInt(c.Param("ccId"), 10, 64)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		pkg.FailCode(c, 400, "请提供日期范围")
		return
	}

	history, err := h.svc.GetCCAttendanceHistory(ccID, startDate, endDate)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, history)
}

// Export 导出在班记录
func (h *AttendanceHandler) Export(c *gin.Context) {
	var query model.CCQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	// 获取日期列表
	dates := c.QueryArray("dates")
	if len(dates) == 0 {
		pkg.FailCode(c, 400, "请提供日期参数")
		return
	}

	fileBytes, err := h.svc.ExportAttendance(&query, dates)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=attendance.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}
