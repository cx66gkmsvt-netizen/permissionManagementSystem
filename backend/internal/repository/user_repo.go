package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: DB}
}

// FindByUserName 根据用户名查找
func (r *UserRepository) FindByUserName(userName string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("user_name = ? AND del_flag = '0'", userName).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找
func (r *UserRepository) FindByID(userID int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Preload("Roles").Preload("Dept").
		Where("user_id = ? AND del_flag = '0'", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 用户列表
func (r *UserRepository) List(query *model.UserQuery, dataScope string) (*model.PageResult, error) {
	var users []model.SysUser
	var total int64

	db := r.db.Model(&model.SysUser{}).Where("del_flag = '0'")

	if query.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+query.UserName+"%")
	}
	if query.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+query.Phone+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.DeptID != nil {
		db = db.Where("dept_id = ?", *query.DeptID)
	}

	// 数据权限过滤
	if dataScope != "" {
		db = db.Where(dataScope)
	}

	db.Count(&total)

	err := db.Preload("Dept").Preload("Roles").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("user_id DESC").
		Find(&users).Error

	return &model.PageResult{Total: total, Rows: users}, err
}

// Create 创建用户
func (r *UserRepository) Create(user *model.SysUser) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.SysUser) error {
	return r.db.Model(user).Updates(user).Error
}

// Delete 删除用户(软删除)
func (r *UserRepository) Delete(userID int64) error {
	return r.db.Model(&model.SysUser{}).Where("user_id = ?", userID).Update("del_flag", "2").Error
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepository) UpdateLoginInfo(userID int64, ip string) error {
	return r.db.Model(&model.SysUser{}).Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"login_ip":   ip,
			"login_date": gorm.Expr("NOW()"),
		}).Error
}

// CheckUserNameUnique 检查用户名唯一性
func (r *UserRepository) CheckUserNameUnique(userName string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.SysUser{}).Where("user_name = ? AND del_flag = '0'", userName)
	if excludeID > 0 {
		query = query.Where("user_id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// GetUserRoles 获取用户角色
func (r *UserRepository) GetUserRoles(userID int64) ([]model.SysRole, error) {
	var user model.SysUser
	err := r.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// SetUserRoles 设置用户角色
func (r *UserRepository) SetUserRoles(userID int64, roleIDs []int64) error {
	// 先删除原有关联
	if err := r.db.Exec("DELETE FROM sys_user_role WHERE user_id = ?", userID).Error; err != nil {
		return err
	}
	// 添加新关联
	for _, roleID := range roleIDs {
		if err := r.db.Exec("INSERT INTO sys_user_role (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
			return err
		}
	}
	return nil
}
