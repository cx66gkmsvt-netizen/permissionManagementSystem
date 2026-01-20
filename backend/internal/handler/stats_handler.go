package handler

import (
	"github.com/gin-gonic/gin"

	"user-center/internal/pkg"
	"user-center/internal/repository"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

// GetStats 获取系统统计数据
func (h *StatsHandler) GetStats(c *gin.Context) {
	db := repository.GetDB()

	var userCount, roleCount, deptCount, menuCount int64

	db.Table("sys_user").Where("del_flag = '0'").Count(&userCount)
	db.Table("sys_role").Where("del_flag = '0'").Count(&roleCount)
	db.Table("sys_dept").Where("del_flag = '0'").Count(&deptCount)
	db.Table("sys_menu").Count(&menuCount)

	pkg.OK(c, gin.H{
		"userCount": userCount,
		"roleCount": roleCount,
		"deptCount": deptCount,
		"menuCount": menuCount,
	})
}
