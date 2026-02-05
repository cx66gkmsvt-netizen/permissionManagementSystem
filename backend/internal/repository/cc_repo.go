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

	// 精确匹配CCID
	if query.CCID != nil {
		db = db.Where("id = ?", *query.CCID)
	}
	// 模糊匹配姓名
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	// 模糊匹配昵称
	if query.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+query.NickName+"%")
	}
	// 精确匹配手机号
	if query.Mobile != "" {
		db = db.Where("mobile = ?", query.Mobile)
	}
	// 角色类型
	if query.RoleType != "" {
		db = db.Where("role_type = ?", query.RoleType)
	}
	// 状态
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	// 是否屏蔽
	if query.IsBlocked != "" {
		db = db.Where("is_blocked = ?", query.IsBlocked)
	}
	// 在班状态
	if query.AttendanceStatus != "" {
		db = db.Where("attendance_status = ?", query.AttendanceStatus)
	}
	// 军团筛选
	if query.LegionID != nil {
		db = db.Where("legion_id = ?", *query.LegionID)
	}
	// 团队筛选
	if query.TeamID != nil {
		db = db.Where("team_id = ?", *query.TeamID)
	}
	// 战队筛选
	if query.SquadID != nil {
		db = db.Where("squad_id = ?", *query.SquadID)
	}

	db.Count(&total)

	err := db.Select("cc_member.*, l.legion_name, t.team_name, s.squad_name").
		Joins("LEFT JOIN cc_legion l ON cc_member.legion_id = l.id").
		Joins("LEFT JOIN cc_team t ON cc_member.team_id = t.id").
		Joins("LEFT JOIN cc_squad s ON cc_member.squad_id = s.id").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("cc_member.performance_rank ASC, cc_member.create_time DESC").
		Find(&list).Error

	// 转换金额为元
	for i := range list {
		list[i].BalanceYuan = float64(list[i].Balance) / 100
		list[i].PerformanceYuan = float64(list[i].MonthlyPerformance) / 100
	}

	return &model.PageResult{Total: total, Rows: list}, err
}

// Get 获取CC信息
func (r *CCRepository) Get(id int64) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.First(&cc, id).Error
	if err == nil {
		cc.BalanceYuan = float64(cc.Balance) / 100
		cc.PerformanceYuan = float64(cc.MonthlyPerformance) / 100
	}
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
	if mobile == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("mobile = ? AND del_flag = '0'", mobile)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckRonglianSeatUnique 检查容联座席号唯一
func (r *CCRepository) CheckRonglianSeatUnique(seat string, excludeID int64) bool {
	if seat == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("ronglian_seat = ? AND del_flag = '0'", seat)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckDiankongSeatUnique 检查点控云座席号唯一
func (r *CCRepository) CheckDiankongSeatUnique(seat string, excludeID int64) bool {
	if seat == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("diankong_seat = ? AND del_flag = '0'", seat)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckHeliAccountUnique 检查合力亿捷账号唯一
func (r *CCRepository) CheckHeliAccountUnique(account string, excludeID int64) bool {
	if account == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("heli_account = ? AND del_flag = '0'", account)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckBaichuanSeatUnique 检查百川智通座席号唯一
func (r *CCRepository) CheckBaichuanSeatUnique(seat string, excludeID int64) bool {
	if seat == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("baichuan_seat = ? AND del_flag = '0'", seat)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// CheckCloudAccount1Unique 检查云客账号1唯一
func (r *CCRepository) CheckCloudAccount1Unique(account string, excludeID int64) bool {
	if account == "" {
		return true
	}
	var count int64
	query := r.db.Model(&model.CCMember{}).Where("cloud_account1 = ? AND del_flag = '0'", account)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// GetByCloudAccount1 根据云客账号1获取CC
func (r *CCRepository) GetByCloudAccount1(account string) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.Where("cloud_account1 = ? AND del_flag = '0'", account).First(&cc).Error
	return &cc, err
}

// GetByDiankongSeat 根据点控云座席号获取CC
func (r *CCRepository) GetByDiankongSeat(seat string) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.Where("diankong_seat = ? AND del_flag = '0'", seat).First(&cc).Error
	return &cc, err
}

// GetByHeliAccount 根据合力亿捷账号获取CC
func (r *CCRepository) GetByHeliAccount(account string) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.Where("heli_account = ? AND del_flag = '0'", account).First(&cc).Error
	return &cc, err
}

// GetByBaichuanSeat 根据百川智通座席号获取CC
func (r *CCRepository) GetByBaichuanSeat(seat string) (*model.CCMember, error) {
	var cc model.CCMember
	err := r.db.Where("baichuan_seat = ? AND del_flag = '0'", seat).First(&cc).Error
	return &cc, err
}

// ListActiveCC 获取在职CC列表（未屏蔽）
func (r *CCRepository) ListActiveCC() ([]model.CCMember, error) {
	var list []model.CCMember
	err := r.db.Where("del_flag = '0' AND is_blocked = '0'").Find(&list).Error
	return list, err
}

// ListBySquadID 根据战队ID获取CC列表
func (r *CCRepository) ListBySquadID(squadID int64) ([]model.CCMember, error) {
	var list []model.CCMember
	err := r.db.Where("squad_id = ? AND del_flag = '0'", squadID).Find(&list).Error
	return list, err
}

// ListByTeamID 根据团队ID获取CC列表
func (r *CCRepository) ListByTeamID(teamID int64) ([]model.CCMember, error) {
	var list []model.CCMember
	err := r.db.Where("team_id = ? AND del_flag = '0'", teamID).Find(&list).Error
	return list, err
}

// ListByLegionID 根据军团ID获取CC列表
func (r *CCRepository) ListByLegionID(legionID int64) ([]model.CCMember, error) {
	var list []model.CCMember
	err := r.db.Where("legion_id = ? AND del_flag = '0'", legionID).Find(&list).Error
	return list, err
}

// UpdateBalance 更新余额
func (r *CCRepository) UpdateBalance(id int64, balance int64) error {
	return r.db.Model(&model.CCMember{}).Where("id = ?", id).Update("balance", balance).Error
}

// UpdateAttendanceStatus 更新在班状态
func (r *CCRepository) UpdateAttendanceStatus(id int64, status string) error {
	return r.db.Model(&model.CCMember{}).Where("id = ?", id).Update("attendance_status", status).Error
}

// ResetAllAttendanceStatus 重置所有CC在班状态为休班
func (r *CCRepository) ResetAllAttendanceStatus() error {
	return r.db.Model(&model.CCMember{}).Where("del_flag = '0'").Update("attendance_status", model.AttendanceStatusOffDuty).Error
}
