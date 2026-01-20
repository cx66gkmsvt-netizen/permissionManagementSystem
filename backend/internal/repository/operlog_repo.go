package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type OperLogRepository struct {
	db *gorm.DB
}

func NewOperLogRepository() *OperLogRepository {
	return &OperLogRepository{db: DB}
}

// Create 创建日志
func (r *OperLogRepository) Create(log *model.SysOperLog) error {
	return r.db.Create(log).Error
}

// List 日志列表
func (r *OperLogRepository) List(pageNum, pageSize int) (*model.PageResult, error) {
	var logs []model.SysOperLog
	var total int64

	r.db.Model(&model.SysOperLog{}).Count(&total)

	offset := (pageNum - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("oper_id DESC").Find(&logs).Error

	return &model.PageResult{Total: total, Rows: logs}, err
}
