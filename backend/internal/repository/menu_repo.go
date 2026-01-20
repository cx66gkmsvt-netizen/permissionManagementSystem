package repository

import (
	"fmt"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{db: DB}
}

// FindByID 根据ID查找
func (r *MenuRepository) FindByID(menuID int64) (*model.SysMenu, error) {
	var menu model.SysMenu
	err := r.db.Where("menu_id = ?", menuID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// SelectAll 查询所有菜单
func (r *MenuRepository) SelectAll() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Where("status = '0'").Order("sort ASC").Find(&menus).Error
	return menus, err
}

// SelectTree 获取菜单树
func (r *MenuRepository) SelectTree() ([]*model.SysMenu, error) {
	menus, err := r.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.SysMenu, parentID int64) []*model.SysMenu {
	var tree []*model.SysMenu
	for i := range menus {
		if menus[i].ParentID == parentID {
			menu := &menus[i]
			menu.Children = buildMenuTree(menus, menu.MenuID)
			tree = append(tree, menu)
		}
	}
	return tree
}

// SelectMenusByRoleIDs 根据角色ID列表查询菜单
func (r *MenuRepository) SelectMenusByRoleIDs(roleIDs []int64) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Distinct().
		Joins("JOIN sys_role_menu rm ON sys_menu.menu_id = rm.menu_id").
		Where("rm.role_id IN ? AND sys_menu.status = '0'", roleIDs).
		Order("sort ASC").
		Find(&menus).Error
	return menus, err
}

// SelectPermsByRoleIDs 根据角色ID列表查询权限标识
func (r *MenuRepository) SelectPermsByRoleIDs(roleIDs []int64) ([]string, error) {
	var perms []string
	err := r.db.Model(&model.SysMenu{}).
		Distinct().
		Joins("JOIN sys_role_menu rm ON sys_menu.menu_id = rm.menu_id").
		Where("rm.role_id IN ? AND sys_menu.status = '0' AND sys_menu.perms != ''", roleIDs).
		Pluck("perms", &perms).Error
	return perms, err
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.SysMenu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.SysMenu) error {
	return r.db.Model(menu).Updates(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(menuID int64) error {
	// 检查是否有子菜单
	var count int64
	r.db.Model(&model.SysMenu{}).Where("parent_id = ?", menuID).Count(&count)
	if count > 0 {
		return fmt.Errorf("存在子菜单，不能删除")
	}
	return r.db.Delete(&model.SysMenu{}, menuID).Error
}
