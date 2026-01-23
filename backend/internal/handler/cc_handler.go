package handler

import (
	"strconv"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/service"

	"github.com/gin-gonic/gin"
)

type CCHandler struct {
	svc *service.CCService
}

func NewCCHandler() *CCHandler {
	return &CCHandler{
		svc: service.NewCCService(),
	}
}

func (h *CCHandler) List(c *gin.Context) {
	var query model.CCQuery
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
