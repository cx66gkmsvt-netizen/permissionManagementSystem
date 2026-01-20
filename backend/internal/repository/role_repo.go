package repository

import (
	"user-center/internal/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: DB}
}

// FindByID 根据ID查找
func (r *RoleRepository) FindByID(roleID int64) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.Preload("Menus").Preload("Depts").
		Where("role_id = ? AND del_flag = '0'", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByKey 根据RoleKey查找
func (r *RoleRepository) FindByKey(roleKey string) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.Where("role_key = ? AND del_flag = '0'", roleKey).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 角色列表
func (r *RoleRepository) List(query *model.RoleQuery) (*model.PageResult, error) {
	var roles []model.SysRole
	var total int64

	db := r.db.Model(&model.SysRole{}).Where("del_flag = '0'")

	if query.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+query.RoleName+"%")
	}
	if query.RoleKey != "" {
		db = db.Where("role_key LIKE ?", "%"+query.RoleKey+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	db.Count(&total)

	err := db.Offset(query.GetOffset()).
		Limit(query.PageSize).
		Order("sort ASC, role_id ASC").
		Find(&roles).Error

	return &model.PageResult{Total: total, Rows: roles}, err
}

// SelectAll 查询所有角色
func (r *RoleRepository) SelectAll() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Where("del_flag = '0' AND status = '0'").Order("sort ASC").Find(&roles).Error
	return roles, err
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.SysRole) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.SysRole) error {
	return r.db.Model(role).Updates(role).Error
}

// Delete 删除角色(软删除)
func (r *RoleRepository) Delete(roleID int64) error {
	return r.db.Model(&model.SysRole{}).Where("role_id = ?", roleID).Update("del_flag", "2").Error
}

// CheckRoleKeyUnique 检查RoleKey唯一性
func (r *RoleRepository) CheckRoleKeyUnique(roleKey string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.SysRole{}).Where("role_key = ? AND del_flag = '0'", roleKey)
	if excludeID > 0 {
		query = query.Where("role_id != ?", excludeID)
	}
	query.Count(&count)
	return count == 0
}

// SetRoleMenus 设置角色菜单
func (r *RoleRepository) SetRoleMenus(roleID int64, menuIDs []int64) error {
	if err := r.db.Exec("DELETE FROM sys_role_menu WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		if err := r.db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", roleID, menuID).Error; err != nil {
			return err
		}
	}
	return nil
}

// SetRoleDepts 设置角色数据权限部门
func (r *RoleRepository) SetRoleDepts(roleID int64, deptIDs []int64) error {
	if err := r.db.Exec("DELETE FROM sys_role_dept WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}
	for _, deptID := range deptIDs {
		if err := r.db.Exec("INSERT INTO sys_role_dept (role_id, dept_id) VALUES (?, ?)", roleID, deptID).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetRoleMenuIDs 获取角色菜单ID
func (r *RoleRepository) GetRoleMenuIDs(roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := r.db.Table("sys_role_menu").Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

// GetRoleDeptIDs 获取角色部门ID
func (r *RoleRepository) GetRoleDeptIDs(roleID int64) ([]int64, error) {
	var deptIDs []int64
	err := r.db.Table("sys_role_dept").Where("role_id = ?", roleID).Pluck("dept_id", &deptIDs).Error
	return deptIDs, err
}
