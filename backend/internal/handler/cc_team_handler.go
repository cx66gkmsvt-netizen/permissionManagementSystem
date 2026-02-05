package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type CCTeamHandler struct {
	svc     *service.CCTeamService
	fundSvc *service.CCFundService
}

func NewCCTeamHandler() *CCTeamHandler {
	return &CCTeamHandler{
		svc:     service.NewCCTeamService(),
		fundSvc: service.NewCCFundService(),
	}
}

// List 获取团队列表
func (h *CCTeamHandler) List(c *gin.Context) {
	var query model.TeamQuery
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

// ListAll 获取全部团队
func (h *CCTeamHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, list)
}

// Get 获取团队详情
func (h *CCTeamHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	team, err := h.svc.Get(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, team)
}

// Create 创建团队
func (h *CCTeamHandler) Create(c *gin.Context) {
	var dto model.TeamCreateDTO
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

// Update 更新团队
func (h *CCTeamHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.TeamUpdateDTO
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

// GetLogs 获取团队修改记录
func (h *CCTeamHandler) GetLogs(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	logs, err := h.svc.GetLogs(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, logs)
}

// GetFund 获取团队资金信息
func (h *CCTeamHandler) GetFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fund, err := h.fundSvc.GetTeamFund(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, fund)
}

// EditFund 编辑团队余额
func (h *CCTeamHandler) EditFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.EditTeamBalance(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Recharge 团队充值
func (h *CCTeamHandler) Recharge(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundRechargeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.TeamRecharge(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Transfer 团队转账
func (h *CCTeamHandler) Transfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundTransferDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.TeamTransfer(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetBills 获取团队账单明细
func (h *CCTeamHandler) GetBills(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	billType := c.DefaultQuery("billType", "non_flow")
	bills, err := h.fundSvc.GetTeamBills(id, billType)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, bills)
}
