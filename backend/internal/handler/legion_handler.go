package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type LegionHandler struct {
	svc     *service.LegionService
	fundSvc *service.CCFundService
}

func NewLegionHandler() *LegionHandler {
	return &LegionHandler{
		svc:     service.NewLegionService(),
		fundSvc: service.NewCCFundService(),
	}
}

// List 获取军团列表
func (h *LegionHandler) List(c *gin.Context) {
	var query model.LegionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	result, err := h.svc.List(&query)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, result)
}

// ListAll 获取全部军团
func (h *LegionHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// Get 获取军团详情
func (h *LegionHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	legion, err := h.svc.Get(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, legion)
}

// Create 创建军团
func (h *LegionHandler) Create(c *gin.Context) {
	var dto model.LegionCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.svc.Create(&dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Update 更新军团
func (h *LegionHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.LegionUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.svc.Update(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetLogs 获取军团跟进记录
func (h *LegionHandler) GetLogs(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	logs, err := h.svc.GetLogs(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, logs)
}

// GetFund 获取军团资金信息
func (h *LegionHandler) GetFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fund, err := h.fundSvc.GetLegionFund(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, fund)
}

// EditFund 编辑军团余额
func (h *LegionHandler) EditFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.EditLegionBalance(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Recharge 军团充值
func (h *LegionHandler) Recharge(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundRechargeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.LegionRecharge(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Transfer 军团转账
func (h *LegionHandler) Transfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundTransferDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.LegionTransfer(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetBills 获取军团账单明细
func (h *LegionHandler) GetBills(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	billType := c.DefaultQuery("billType", "non_flow")
	bills, err := h.fundSvc.GetLegionBills(id, billType)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, bills)
}

// getOperatorInfo 获取操作人信息
func getOperatorInfo(c *gin.Context) (int64, string) {
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	var operatorID int64
	var operatorName string

	if id, ok := userID.(int64); ok {
		operatorID = id
	}
	if name, ok := userName.(string); ok {
		operatorName = name
	}

	return operatorID, operatorName
}
