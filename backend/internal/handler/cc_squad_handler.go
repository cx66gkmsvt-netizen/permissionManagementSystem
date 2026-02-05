package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type CCSquadHandler struct {
	svc     *service.CCSquadService
	fundSvc *service.CCFundService
}

func NewCCSquadHandler() *CCSquadHandler {
	return &CCSquadHandler{
		svc:     service.NewCCSquadService(),
		fundSvc: service.NewCCFundService(),
	}
}

// List 获取战队列表
func (h *CCSquadHandler) List(c *gin.Context) {
	var query model.SquadQuery
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

// ListAll 获取全部战队
func (h *CCSquadHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// ListByTeam 根据团队获取战队列表
func (h *CCSquadHandler) ListByTeam(c *gin.Context) {
	teamID, _ := strconv.ParseInt(c.Query("teamId"), 10, 64)
	list, err := h.svc.ListByTeamID(teamID)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// Get 获取战队详情
func (h *CCSquadHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	squad, err := h.svc.Get(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, squad)
}

// Create 创建战队
func (h *CCSquadHandler) Create(c *gin.Context) {
	var dto model.SquadCreateDTO
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

// Update 更新战队
func (h *CCSquadHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.SquadUpdateDTO
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

// GetLogs 获取战队修改记录
func (h *CCSquadHandler) GetLogs(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	logs, err := h.svc.GetLogs(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, logs)
}

// GetFund 获取战队资金信息
func (h *CCSquadHandler) GetFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fund, err := h.fundSvc.GetSquadFund(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, fund)
}

// EditFund 编辑战队余额
func (h *CCSquadHandler) EditFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.EditSquadBalance(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Recharge 战队充值
func (h *CCSquadHandler) Recharge(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundRechargeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.SquadRecharge(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Transfer 战队转账
func (h *CCSquadHandler) Transfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundTransferDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.SquadTransfer(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetBills 获取战队账单明细
func (h *CCSquadHandler) GetBills(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	billType := c.DefaultQuery("billType", "non_flow")
	bills, err := h.fundSvc.GetSquadBills(id, billType)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, bills)
}
