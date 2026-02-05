package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type CCSquadRepository struct {
	db *gorm.DB
}

func NewCCSquadRepository() *CCSquadRepository {
	return &CCSquadRepository{db: DB}
}

// List 战队列表
func (r *CCSquadRepository) List(query *model.SquadQuery) (*model.PageResult, error) {
	var list []model.CCSquad
	var total int64

	db := r.db.Model(&model.CCSquad{}).Where("cc_squad.del_flag = '0'")

	if query.SquadName != "" {
		db = db.Where("cc_squad.squad_name LIKE ?", "%"+query.SquadName+"%")
	}
	if query.TeamID != nil {
		db = db.Where("cc_squad.team_id = ?", *query.TeamID)
	}

	db.Count(&total)

	err := db.Select("cc_squad.*, m.name as leader_name, t.team_name, tm.name as team_leader_name, l.legion_name, lm.name as legion_leader_name").
		Joins("LEFT JOIN cc_member m ON cc_squad.leader_id = m.id").
		Joins("LEFT JOIN cc_team t ON cc_squad.team_id = t.id").
		Joins("LEFT JOIN cc_member tm ON t.leader_id = tm.id").
		Joins("LEFT JOIN cc_legion l ON t.legion_id = l.id").
		Joins("LEFT JOIN cc_member lm ON l.leader_id = lm.id").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("cc_squad.id ASC").
		Find(&list).Error

	return &model.PageResult{Total: total, Rows: list}, err
}

// ListAll 获取全部战队（不分页）
func (r *CCSquadRepository) ListAll() ([]model.CCSquad, error) {
	var list []model.CCSquad
	err := r.db.Where("del_flag = '0'").Order("id ASC").Find(&list).Error
	return list, err
}

// ListByTeamID 根据团队ID获取战队列表
func (r *CCSquadRepository) ListByTeamID(teamID int64) ([]model.CCSquad, error) {
	var list []model.CCSquad
	err := r.db.Where("team_id = ? AND del_flag = '0'", teamID).Find(&list).Error
	return list, err
}

// Get 获取战队信息
func (r *CCSquadRepository) Get(id int64) (*model.CCSquad, error) {
	var squad model.CCSquad
	err := r.db.First(&squad, id).Error
	return &squad, err
}

// GetByName 根据名称获取战队
func (r *CCSquadRepository) GetByName(name string) (*model.CCSquad, error) {
	var squad model.CCSquad
	err := r.db.Where("squad_name = ? AND del_flag = '0'", name).First(&squad).Error
	return &squad, err
}

// Create 创建战队
func (r *CCSquadRepository) Create(squad *model.CCSquad) error {
	return r.db.Create(squad).Error
}

// Update 更新战队
func (r *CCSquadRepository) Update(squad *model.CCSquad) error {
	return r.db.Model(squad).Updates(squad).Error
}

// Delete 删除战队
func (r *CCSquadRepository) Delete(id int64) error {
	return r.db.Model(&model.CCSquad{}).Where("id = ?", id).Update("del_flag", "2").Error
}

// CheckNameUnique 检查战队名称唯一
func (r *CCSquadRepository) CheckNameUnique(name string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.CCSquad{}).Where("squad_name = ? AND del_flag = '0'", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// UpdateBalance 更新余额
func (r *CCSquadRepository) UpdateBalance(id int64, balance int64) error {
	return r.db.Model(&model.CCSquad{}).Where("id = ?", id).Update("balance", balance).Error
}

// UpdateLeader 更新战队长
func (r *CCSquadRepository) UpdateLeader(id int64, leaderID *int64) error {
	return r.db.Model(&model.CCSquad{}).Where("id = ?", id).Update("leader_id", leaderID).Error
}

// UpdateTeam 更新所属团队
func (r *CCSquadRepository) UpdateTeam(id int64, teamID int64) error {
	return r.db.Model(&model.CCSquad{}).Where("id = ?", id).Update("team_id", teamID).Error
}

// GetBalance 获取战队余额
func (r *CCSquadRepository) GetBalance(id int64) (int64, error) {
	var squad model.CCSquad
	err := r.db.Select("balance").First(&squad, id).Error
	return squad.Balance, err
}

// CountByTeamID 统计团队下的战队数量
func (r *CCSquadRepository) CountByTeamID(teamID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.CCSquad{}).Where("team_id = ? AND del_flag = '0'", teamID).Count(&count).Error
	return count, err
}

// GetSquadLeaderID 获取战队长ID
func (r *CCSquadRepository) GetSquadLeaderID(squadID int64) (*int64, error) {
	var squad model.CCSquad
	err := r.db.Select("leader_id").First(&squad, squadID).Error
	return squad.LeaderID, err
}
