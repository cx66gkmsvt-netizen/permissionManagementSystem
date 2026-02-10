package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type CCHandler struct {
	svc     *service.CCService
	fundSvc *service.CCFundService
}

func NewCCHandler() *CCHandler {
	return &CCHandler{
		svc:     service.NewCCService(),
		fundSvc: service.NewCCFundService(),
	}
}

func (h *CCHandler) List(c *gin.Context) {
	var query model.CCQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	// 根据当前用户角色应用数据范围过滤
	scope := getDataScope(c)
	applyDataScopeToCCQuery(&query, scope)

	result, err := h.svc.List(&query)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, result)
}

func (h *CCHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cc, err := h.svc.Get(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, cc)
}

func (h *CCHandler) Create(c *gin.Context) {
	var cc model.CCMember
	if err := c.ShouldBindJSON(&cc); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	if err := h.svc.Create(&cc); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

func (h *CCHandler) Update(c *gin.Context) {
	var cc model.CCMember
	if err := c.ShouldBindJSON(&cc); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}
	cc.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.svc.Update(&cc); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

func (h *CCHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetFund 获取CC资金信息
func (h *CCHandler) GetFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fund, err := h.fundSvc.GetCCFund(id)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, fund)
}

// EditFund 编辑CC余额
func (h *CCHandler) EditFund(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundEditDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.EditCCBalance(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// Transfer CC转账
func (h *CCHandler) Transfer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto model.FundTransferDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		pkg.FailCode(c, 400, err.Error())
		return
	}

	operatorID, operatorName := getOperatorInfo(c)
	if err := h.fundSvc.CCTransfer(id, &dto, operatorID, operatorName); err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, nil)
}

// GetBills 获取CC账单明细
func (h *CCHandler) GetBills(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	billType := c.DefaultQuery("billType", "non_flow")
	bills, err := h.fundSvc.GetCCBills(id, billType)
	if err != nil {
		pkg.Fail(c, err.Error())
		return
	}
	pkg.OK(c, bills)
}
