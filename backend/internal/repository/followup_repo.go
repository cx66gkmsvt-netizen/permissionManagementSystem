package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type FollowUpRepository struct {
	db *gorm.DB
}

func NewFollowUpRepository() *FollowUpRepository {
	return &FollowUpRepository{db: DB}
}

// Create 创建跟进记录
func (r *FollowUpRepository) Create(record *model.SysFollowUp) error {
	return r.db.Create(record).Error
}

// ListByTarget 根据目标查询跟进记录
func (r *FollowUpRepository) ListByTarget(targetType string, targetID int64) ([]model.SysFollowUp, error) {
	var list []model.SysFollowUp
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("oper_time DESC").
		Find(&list).Error
	return list, err
}
