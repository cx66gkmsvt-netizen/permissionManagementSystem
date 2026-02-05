package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type LegionRepository struct {
	db *gorm.DB
}

func NewLegionRepository() *LegionRepository {
	return &LegionRepository{db: DB}
}

// List 军团列表
func (r *LegionRepository) List(query *model.LegionQuery) (*model.PageResult, error) {
	var list []model.CCLegion
	var total int64

	db := r.db.Model(&model.CCLegion{}).Where("del_flag = '0'")

	if query.LegionName != "" {
		db = db.Where("legion_name LIKE ?", "%"+query.LegionName+"%")
	}

	db.Count(&total)

	err := db.Select("cc_legion.*, m.name as leader_name").
		Joins("LEFT JOIN cc_member m ON cc_legion.leader_id = m.id").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("cc_legion.id ASC").
		Find(&list).Error

	return &model.PageResult{Total: total, Rows: list}, err
}

// ListAll 获取全部军团（不分页）
func (r *LegionRepository) ListAll() ([]model.CCLegion, error) {
	var list []model.CCLegion
	err := r.db.Where("del_flag = '0'").Order("id ASC").Find(&list).Error
	return list, err
}

// Get 获取军团信息
func (r *LegionRepository) Get(id int64) (*model.CCLegion, error) {
	var legion model.CCLegion
	err := r.db.First(&legion, id).Error
	return &legion, err
}

// GetByName 根据名称获取军团
func (r *LegionRepository) GetByName(name string) (*model.CCLegion, error) {
	var legion model.CCLegion
	err := r.db.Where("legion_name = ? AND del_flag = '0'", name).First(&legion).Error
	return &legion, err
}

// Create 创建军团
func (r *LegionRepository) Create(legion *model.CCLegion) error {
	return r.db.Create(legion).Error
}

// Update 更新军团
func (r *LegionRepository) Update(legion *model.CCLegion) error {
	return r.db.Model(legion).Updates(legion).Error
}

// Delete 删除军团
func (r *LegionRepository) Delete(id int64) error {
	return r.db.Model(&model.CCLegion{}).Where("id = ?", id).Update("del_flag", "2").Error
}

// CheckNameUnique 检查军团名称唯一
func (r *LegionRepository) CheckNameUnique(name string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.CCLegion{}).Where("legion_name = ? AND del_flag = '0'", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// UpdateBalance 更新余额
func (r *LegionRepository) UpdateBalance(id int64, balance int64) error {
	return r.db.Model(&model.CCLegion{}).Where("id = ?", id).Update("balance", balance).Error
}

// UpdateLeader 更新军团长
func (r *LegionRepository) UpdateLeader(id int64, leaderID *int64) error {
	return r.db.Model(&model.CCLegion{}).Where("id = ?", id).Update("leader_id", leaderID).Error
}

// GetBalance 获取军团余额
func (r *LegionRepository) GetBalance(id int64) (int64, error) {
	var legion model.CCLegion
	err := r.db.Select("balance").First(&legion, id).Error
	return legion.Balance, err
}
