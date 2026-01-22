package repository

import (
	"fmt"
	"time"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type OperLogRepository struct {
	db *gorm.DB
}

func NewOperLogRepository() *OperLogRepository {
	return &OperLogRepository{db: DB}
}

// getTableName 获取当前月份的表名
func (r *OperLogRepository) getTableName() string {
	return fmt.Sprintf("sys_oper_log_%s", time.Now().Format("200601"))
}

// Create 创建日志
func (r *OperLogRepository) Create(log *model.SysOperLog) error {
	tableName := r.getTableName()

	// 检查表是否存在，不存在则创建
	if !r.db.Migrator().HasTable(tableName) {
		// 使用 LIKE 语法复制表结构，包含索引和注释
		err := r.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE sys_oper_log", tableName)).Error
		if err != nil {
			return err
		}
	}

	return r.db.Table(tableName).Create(log).Error
}

// List 日志列表 (默认查询当月数据)
func (r *OperLogRepository) List(pageNum, pageSize int) (*model.PageResult, error) {
	var logs []model.SysOperLog
	var total int64
	tableName := r.getTableName()

	// 检查表是否存在
	if !r.db.Migrator().HasTable(tableName) {
		return &model.PageResult{Total: 0, Rows: []model.SysOperLog{}}, nil
	}

	r.db.Table(tableName).Count(&total)

	offset := (pageNum - 1) * pageSize
	err := r.db.Table(tableName).Offset(offset).Limit(pageSize).Order("oper_id DESC").Find(&logs).Error

	return &model.PageResult{Total: total, Rows: logs}, err
}
