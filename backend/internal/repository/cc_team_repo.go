package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type CCTeamRepository struct {
	db *gorm.DB
}

func NewCCTeamRepository() *CCTeamRepository {
	return &CCTeamRepository{db: DB}
}

// List 团队列表
func (r *CCTeamRepository) List(query *model.TeamQuery) (*model.PageResult, error) {
	var list []model.CCTeam
	var total int64

	db := r.db.Model(&model.CCTeam{}).Where("del_flag = '0'")

	if query.TeamName != "" {
		db = db.Where("team_name LIKE ?", "%"+query.TeamName+"%")
	}
	if query.BusinessType != "" {
		db = db.Where("business_type = ?", query.BusinessType)
	}
	if query.LegionID != nil {
		db = db.Where("legion_id = ?", *query.LegionID)
	}

	db.Count(&total)

	err := db.Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("id ASC").
		Find(&list).Error

	return &model.PageResult{Total: total, Rows: list}, err
}

// ListAll 获取全部团队（不分页）
func (r *CCTeamRepository) ListAll() ([]model.CCTeam, error) {
	var list []model.CCTeam
	err := r.db.Where("del_flag = '0'").Order("id ASC").Find(&list).Error
	return list, err
}

// ListByLegionID 根据军团ID获取团队列表
func (r *CCTeamRepository) ListByLegionID(legionID int64) ([]model.CCTeam, error) {
	var list []model.CCTeam
	err := r.db.Where("legion_id = ? AND del_flag = '0'", legionID).Find(&list).Error
	return list, err
}

// Get 获取团队信息
func (r *CCTeamRepository) Get(id int64) (*model.CCTeam, error) {
	var team model.CCTeam
	err := r.db.First(&team, id).Error
	return &team, err
}

// GetByName 根据名称获取团队
func (r *CCTeamRepository) GetByName(name string) (*model.CCTeam, error) {
	var team model.CCTeam
	err := r.db.Where("team_name = ? AND del_flag = '0'", name).First(&team).Error
	return &team, err
}

// Create 创建团队
func (r *CCTeamRepository) Create(team *model.CCTeam) error {
	return r.db.Create(team).Error
}

// Update 更新团队
func (r *CCTeamRepository) Update(team *model.CCTeam) error {
	return r.db.Model(team).Updates(team).Error
}

// Delete 删除团队
func (r *CCTeamRepository) Delete(id int64) error {
	return r.db.Model(&model.CCTeam{}).Where("id = ?", id).Update("del_flag", "2").Error
}

// CheckNameUnique 检查团队名称唯一
func (r *CCTeamRepository) CheckNameUnique(name string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.CCTeam{}).Where("team_name = ? AND del_flag = '0'", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// UpdateBalance 更新余额
func (r *CCTeamRepository) UpdateBalance(id int64, balance int64) error {
	return r.db.Model(&model.CCTeam{}).Where("id = ?", id).Update("balance", balance).Error
}

// UpdateLeader 更新团长
func (r *CCTeamRepository) UpdateLeader(id int64, leaderID *int64) error {
	return r.db.Model(&model.CCTeam{}).Where("id = ?", id).Update("leader_id", leaderID).Error
}

// UpdateLegion 更新所属军团
func (r *CCTeamRepository) UpdateLegion(id int64, legionID *int64) error {
	return r.db.Model(&model.CCTeam{}).Where("id = ?", id).Update("legion_id", legionID).Error
}

// GetBalance 获取团队余额
func (r *CCTeamRepository) GetBalance(id int64) (int64, error) {
	var team model.CCTeam
	err := r.db.Select("balance").First(&team, id).Error
	return team.Balance, err
}

// CountByLegionID 统计军团下的团队数量
func (r *CCTeamRepository) CountByLegionID(legionID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.CCTeam{}).Where("legion_id = ? AND del_flag = '0'", legionID).Count(&count).Error
	return count, err
}
