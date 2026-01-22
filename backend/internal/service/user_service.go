package service

import (
	"errors"
	"fmt"
	"strings"

	"user-center/internal/model"
	"user-center/internal/repository"
)

type UserService struct {
	userRepo        *repository.UserRepository
	roleRepo        *repository.RoleRepository
	deptRepo        *repository.DeptRepository
	followUpService *FollowUpService
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:        repository.NewUserRepository(),
		roleRepo:        repository.NewRoleRepository(),
		deptRepo:        repository.NewDeptRepository(),
		followUpService: NewFollowUpService(),
	}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID int64) (*model.SysUser, error) {
	return s.userRepo.FindByID(userID)
}

// GetUserByUserName 根据用户名获取用户
func (s *UserService) GetUserByUserName(userName string) (*model.SysUser, error) {
	return s.userRepo.FindByUserName(userName)
}

// List 用户列表(带数据权限)
func (s *UserService) List(query *model.UserQuery, currentUser *model.SysUser) (*model.PageResult, error) {
	dataScope := s.buildDataScope(currentUser)
	return s.userRepo.List(query, dataScope)
}

// Create 创建用户
func (s *UserService) Create(req *model.CreateUserRequest, operatorID int64) error {
	// 检查用户名唯一性
	if !s.userRepo.CheckUserNameUnique(req.UserName, 0) {
		return errors.New("用户名已存在")
	}

	user := &model.SysUser{
		DeptID:   req.DeptID,
		UserName: req.UserName,
		NickName: req.NickName,
		Password: req.Password, // 已在handler加密
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		CreateBy: &operatorID,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// 记录跟进
	s.followUpService.Record("sys_user", user.UserID, fmt.Sprintf("创建用户: %s", user.UserName), operatorID, "")

	// 设置角色
	if len(req.RoleIDs) > 0 {
		return s.userRepo.SetUserRoles(user.UserID, req.RoleIDs)
	}
	return nil
}

// Update 更新用户
func (s *UserService) Update(userID int64, req *model.UpdateUserRequest, operatorID int64) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 禁止修改admin
	if user.UserID == 1 {
		return errors.New("不允许修改超级管理员")
	}

	user.DeptID = req.DeptID
	user.NickName = req.NickName
	user.Email = req.Email
	user.Phone = req.Phone
	user.Status = req.Status

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 记录跟进
	s.followUpService.Record("sys_user", userID, "更新用户资料", operatorID, "")

	// 更新角色
	if req.RoleIDs != nil {
		return s.userRepo.SetUserRoles(userID, req.RoleIDs)
	}
	return nil
}

// Delete 删除用户
func (s *UserService) Delete(userID int64, operatorID int64) error {
	if userID == 1 {
		return errors.New("不允许删除超级管理员")
	}
	if err := s.userRepo.Delete(userID); err != nil {
		return err
	}

	// 记录跟进
	return s.followUpService.Record("sys_user", userID, "删除用户", operatorID, "")
}

// GetUserRoles 获取用户角色
func (s *UserService) GetUserRoles(userID int64) ([]model.SysRole, error) {
	return s.userRepo.GetUserRoles(userID)
}

// buildDataScope 构建数据权限SQL
func (s *UserService) buildDataScope(user *model.SysUser) string {
	// 超级管理员不过滤
	if user.UserID == 1 {
		return ""
	}

	roles, err := s.userRepo.GetUserRoles(user.UserID)
	if err != nil || len(roles) == 0 {
		return "1=0" // 无角色，无权限
	}

	var conditions []string
	for _, role := range roles {
		switch role.DataScope {
		case "1": // 全部数据权限
			return ""
		case "2": // 自定义数据权限
			conditions = append(conditions, fmt.Sprintf(
				"dept_id IN (SELECT dept_id FROM sys_role_dept WHERE role_id = %d)", role.RoleID))
		case "3": // 本部门及以下
			if user.DeptID != nil {
				conditions = append(conditions, fmt.Sprintf(
					"dept_id IN (SELECT dept_id FROM sys_dept WHERE dept_id = %d OR FIND_IN_SET(%d, ancestors))",
					*user.DeptID, *user.DeptID))
			}
		case "4": // 仅本部门
			if user.DeptID != nil {
				conditions = append(conditions, fmt.Sprintf("dept_id = %d", *user.DeptID))
			}
		case "5": // 仅本人数据
			conditions = append(conditions, fmt.Sprintf("create_by = %d", user.UserID))
		}
	}

	if len(conditions) == 0 {
		return "1=0"
	}
	return "(" + strings.Join(conditions, " OR ") + ")"
}
