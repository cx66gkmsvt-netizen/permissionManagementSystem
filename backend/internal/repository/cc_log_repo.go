package repository

import (
	"time"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type CCLogRepository struct {
	db *gorm.DB
}

func NewCCLogRepository() *CCLogRepository {
	return &CCLogRepository{db: DB}
}

// ==================== 管理日志 ====================

// CreateManageLog 创建管理日志
func (r *CCLogRepository) CreateManageLog(log *model.CCManageLog) error {
	log.CreateTime = time.Now()
	return r.db.Create(log).Error
}

// ListManageLogsByTarget 根据目标获取管理日志列表
func (r *CCLogRepository) ListManageLogsByTarget(targetType string, targetID int64) ([]model.CCManageLog, error) {
	var list []model.CCManageLog
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ListManageLogsByLegion 获取军团相关的所有管理日志
func (r *CCLogRepository) ListManageLogsByLegion(legionID int64) ([]model.CCManageLog, error) {
	var list []model.CCManageLog
	err := r.db.Where("(target_type = 'legion' AND target_id = ?) OR "+
		"(target_type = 'team' AND target_id IN (SELECT id FROM cc_team WHERE legion_id = ?)) OR "+
		"(target_type = 'squad' AND target_id IN (SELECT id FROM cc_squad WHERE team_id IN (SELECT id FROM cc_team WHERE legion_id = ?))) OR "+
		"(target_type = 'cc' AND target_id IN (SELECT id FROM cc_member WHERE legion_id = ?))",
		legionID, legionID, legionID, legionID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ListManageLogsByTeam 获取团队相关的所有管理日志
func (r *CCLogRepository) ListManageLogsByTeam(teamID int64) ([]model.CCManageLog, error) {
	var list []model.CCManageLog
	err := r.db.Where("(target_type = 'team' AND target_id = ?) OR "+
		"(target_type = 'squad' AND target_id IN (SELECT id FROM cc_squad WHERE team_id = ?)) OR "+
		"(target_type = 'cc' AND target_id IN (SELECT id FROM cc_member WHERE team_id = ?))",
		teamID, teamID, teamID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ListManageLogsBySquad 获取战队相关的所有管理日志
func (r *CCLogRepository) ListManageLogsBySquad(squadID int64) ([]model.CCManageLog, error) {
	var list []model.CCManageLog
	err := r.db.Where("(target_type = 'squad' AND target_id = ?) OR "+
		"(target_type = 'cc' AND target_id IN (SELECT id FROM cc_member WHERE squad_id = ?))",
		squadID, squadID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ListManageLogsByCC 获取CC相关的所有管理日志
func (r *CCLogRepository) ListManageLogsByCC(ccID int64) ([]model.CCManageLog, error) {
	var list []model.CCManageLog
	err := r.db.Where("target_type = 'cc' AND target_id = ?", ccID).
		Order("create_time DESC").
		Find(&list).Error
	return list, err
}

// ==================== 资金日志 ====================

// CreateFundLog 创建资金日志
func (r *CCLogRepository) CreateFundLog(log *model.CCFundLog) error {
	log.CreateTime = time.Now()
	return r.db.Create(log).Error
}

// ListFundLogs 获取资金日志列表
func (r *CCLogRepository) ListFundLogs(query *model.FundLogQuery) (*model.PageResult, error) {
	var list []model.CCFundLog
	var total int64

	db := r.db.Model(&model.CCFundLog{})

	if query.TargetType != "" {
		db = db.Where("target_type = ?", query.TargetType)
	}
	if query.TargetID > 0 {
		db = db.Where("target_id = ?", query.TargetID)
	}
	if query.LogType != "" {
		db = db.Where("log_type = ?", query.LogType)
	}

	// 账单类型筛选
	if query.BillType == "flow" {
		// 流水：业绩流水增加/减少
		db = db.Where("log_type IN ?", []string{
			model.FundLogTypePerformanceIncrease,
			model.FundLogTypePerformanceDecrease,
		})
	} else if query.BillType == "non_flow" {
		// 非流水：其他类型
		db = db.Where("log_type NOT IN ?", []string{
			model.FundLogTypePerformanceIncrease,
			model.FundLogTypePerformanceDecrease,
		})
	}
	// 默认全部，不加筛选

	db.Count(&total)

	err := db.Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("id DESC").
		Find(&list).Error

	return &model.PageResult{Total: total, Rows: list}, err
}

// ListFundLogsByTarget 根据目标获取资金日志列表（不分页）
func (r *CCLogRepository) ListFundLogsByTarget(targetType string, targetID int64, billType string) ([]model.CCFundLog, error) {
	var list []model.CCFundLog

	db := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID)

	if billType == "flow" {
		db = db.Where("log_type IN ?", []string{
			model.FundLogTypePerformanceIncrease,
			model.FundLogTypePerformanceDecrease,
		})
	} else if billType == "non_flow" {
		db = db.Where("log_type NOT IN ?", []string{
			model.FundLogTypePerformanceIncrease,
			model.FundLogTypePerformanceDecrease,
		})
	}

	err := db.Order("id DESC").Find(&list).Error
	return list, err
}

// GetLatestFundLog 获取最新的资金日志
func (r *CCLogRepository) GetLatestFundLog(targetType string, targetID int64) (*model.CCFundLog, error) {
	var log model.CCFundLog
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("id DESC").
		First(&log).Error
	return &log, err
}
