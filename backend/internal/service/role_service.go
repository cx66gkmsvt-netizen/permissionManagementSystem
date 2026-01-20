package service

import (
	"errors"

	"user-center/internal/model"
	"user-center/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
	}
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(roleID int64) (*model.SysRole, error) {
	return s.roleRepo.FindByID(roleID)
}

// List 角色列表
func (s *RoleService) List(query *model.RoleQuery) (*model.PageResult, error) {
	return s.roleRepo.List(query)
}

// SelectAll 查询所有角色
func (s *RoleService) SelectAll() ([]model.SysRole, error) {
	return s.roleRepo.SelectAll()
}

// Create 创建角色
func (s *RoleService) Create(req *model.CreateRoleRequest) error {
	// 检查RoleKey唯一性
	if !s.roleRepo.CheckRoleKeyUnique(req.RoleKey, 0) {
		return errors.New("角色权限字符已存在")
	}

	role := &model.SysRole{
		RoleName:  req.RoleName,
		RoleKey:   req.RoleKey,
		DataScope: req.DataScope,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	if role.DataScope == "" {
		role.DataScope = "1" // 默认全部数据
	}

	if err := s.roleRepo.Create(role); err != nil {
		return err
	}

	// 设置菜单权限
	if len(req.MenuIDs) > 0 {
		if err := s.roleRepo.SetRoleMenus(role.RoleID, req.MenuIDs); err != nil {
			return err
		}
	}

	// 设置数据权限部门(自定义时)
	if role.DataScope == "2" && len(req.DeptIDs) > 0 {
		return s.roleRepo.SetRoleDepts(role.RoleID, req.DeptIDs)
	}

	return nil
}

// Update 更新角色
func (s *RoleService) Update(roleID int64, req *model.UpdateRoleRequest) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 禁止修改admin角色
	if role.RoleKey == "admin" {
		return errors.New("不允许修改超级管理员角色")
	}

	// 检查RoleKey唯一性
	if req.RoleKey != "" && req.RoleKey != role.RoleKey {
		if !s.roleRepo.CheckRoleKeyUnique(req.RoleKey, roleID) {
			return errors.New("角色权限字符已存在")
		}
		role.RoleKey = req.RoleKey
	}

	role.RoleName = req.RoleName
	role.DataScope = req.DataScope
	role.Sort = req.Sort
	role.Status = req.Status
	role.Remark = req.Remark

	if err := s.roleRepo.Update(role); err != nil {
		return err
	}

	// 更新菜单权限
	if req.MenuIDs != nil {
		if err := s.roleRepo.SetRoleMenus(roleID, req.MenuIDs); err != nil {
			return err
		}
	}

	// 更新数据权限部门
	if role.DataScope == "2" && req.DeptIDs != nil {
		return s.roleRepo.SetRoleDepts(roleID, req.DeptIDs)
	}

	return nil
}

// Delete 删除角色
func (s *RoleService) Delete(roleID int64) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.RoleKey == "admin" {
		return errors.New("不允许删除超级管理员角色")
	}
	return s.roleRepo.Delete(roleID)
}

// GetRoleMenuIDs 获取角色菜单ID
func (s *RoleService) GetRoleMenuIDs(roleID int64) ([]int64, error) {
	return s.roleRepo.GetRoleMenuIDs(roleID)
}

// GetRoleDeptIDs 获取角色部门ID
func (s *RoleService) GetRoleDeptIDs(roleID int64) ([]int64, error) {
	return s.roleRepo.GetRoleDeptIDs(roleID)
}
