package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type CCRepository struct {
	db *gorm.DB
}

func NewCCRepository() *CCRepository {
	return &CCRepository{db: DB}
}

// List CC列表
func (r *CCRepository) List(query *model.CCQuery) (*model.PageResult, error) {
	var list []model.CCMember
	var total int64

	db := r.db.Model(&model.CCMember{}).Where("del_flag = '0'")

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Mobile != "" {
		db = db.Where("mobile LIKE ?", "%"+query.Mobile+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.TeamID != nil {
		db = db.Where("team_id = ?", *query.TeamID)
	}
	if query.SquadID != nil {
		db = db.Where("squad_id = ?", *query.SquadID)
	}

	db.Count(&total)

	err := db.Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("create_time DESC").
		Find(&list).Error

	return &model.PageResult{Total: total, Rows: list}, err
}

// Get 获取CC信息
func (r *CCRepository) Get(id int64) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.First(&cc, id).Error
	return &cc, err
}

// Create 创建CC
func (r *CCRepository) Create(cc *model.CCMember) error {
	return r.db.Create(cc).Error
}

// Update 更新CC
func (r *CCRepository) Update(cc *model.CCMember) error {
	return r.db.Model(cc).Updates(cc).Error
}

// Delete 删除CC
func (r *CCRepository) Delete(id int64) error {
	return r.db.Model(&model.CCMember{}).Where("id = ?", id).Update("del_flag", "2").Error
}

// CheckMobileUnique 检查手机号唯一
func (r *CCRepository) CheckMobileUnique(mobile string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("mobile = ? AND del_flag = '0'", mobile)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckCnoUnique 检查座席号唯一
func (r *CCRepository) CheckCnoUnique(cno string, excludeID int64) bool {
	if cno == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("cno = ? AND del_flag = '0'", cno)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}
