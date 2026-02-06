package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"user-center/internal/model"
	"user-center/internal/pkg/trace"
	"user-center/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	deptRepo *repository.DeptRepository
	ccRepo   *repository.CCRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
		roleRepo: repository.NewRoleRepository(),
		deptRepo: repository.NewDeptRepository(),
		ccRepo:   repository.NewCCRepository(),
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
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest, operatorID int64) error {
	trace.AddStep(ctx, "Start Create User", "UserName: %s", req.UserName)

	// 检查用户名唯一性
	if !s.userRepo.CheckUserNameUnique(req.UserName, 0) {
		trace.AddStep(ctx, "Check Unique Failed", "Username %s already exists", req.UserName)
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

	trace.AddStep(ctx, "DB Create", "Saving user to database")
	if err := s.userRepo.Create(user); err != nil {
		trace.AddStep(ctx, "DB Create Error", "Error: %v", err)
		return err
	}

	// 设置角色
	if len(req.RoleIDs) > 0 {
		trace.AddStep(ctx, "Set Roles", "Assigning roles: %v", req.RoleIDs)
		if err := s.userRepo.SetUserRoles(user.UserID, req.RoleIDs); err != nil {
			return err
		}

		// 同步到CC管理
		if err := s.syncUserToCC(ctx, user, req.RoleIDs); err != nil {
			// 仅记录错误，不阻断主流程
			trace.AddStep(ctx, "Sync CC Failed", "Error: %v", err)
		}
	}
	return nil
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, userID int64, req *model.UpdateUserRequest) error {
	trace.AddStep(ctx, "Start Update User", "UserID: %d", userID)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		trace.AddStep(ctx, "Find User Failed", "User not found")
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

	trace.AddStep(ctx, "DB Update", "Updating user record")
	if err := s.userRepo.Update(user); err != nil {
		trace.AddStep(ctx, "DB Update Failed", "Error: %v", err)
		return err
	}

	// 更新角色
	if req.RoleIDs != nil {
		trace.AddStep(ctx, "Update Roles", "New roles: %v", req.RoleIDs)
		if err := s.userRepo.SetUserRoles(userID, req.RoleIDs); err != nil {
			return err
		}

		// 同步到CC管理
		if err := s.syncUserToCC(ctx, user, req.RoleIDs); err != nil {
			trace.AddStep(ctx, "Sync CC Failed", "Error: %v", err)
		}
	} else {
		// 如果只更新了用户信息，也尝试同步基本信息到CC
		// 需要先获取当前角色
		roles, err := s.userRepo.GetUserRoles(userID)
		if err == nil && len(roles) > 0 {
			var roleIDs []int64
			for _, r := range roles {
				roleIDs = append(roleIDs, r.RoleID)
			}
			if err := s.syncUserToCC(ctx, user, roleIDs); err != nil {
				trace.AddStep(ctx, "Sync CC Failed", "Error: %v", err)
			}
		}
	}
	return nil
}

// syncUserToCC 同步用户到CC表
func (s *UserService) syncUserToCC(ctx context.Context, user *model.SysUser, roleIDs []int64) error {
	// 1. 获取角色Key
	roles, err := s.roleRepo.FindByIDs(roleIDs)
	if err != nil {
		return err
	}

	// 2. 检查是否包含CC相关角色
	var ccRoleKey string
	var isCCUser bool
	for _, role := range roles {
		if role.RoleKey == "cc" {
			ccRoleKey = model.RoleTypeCC
			isCCUser = true
			break
		} else if role.RoleKey == "cc_team_leader" {
			ccRoleKey = model.RoleTypeTeamLeader
			isCCUser = true
			break
		} else if role.RoleKey == "cc_squad_leader" {
			ccRoleKey = model.RoleTypeSquadLeader
			isCCUser = true
			break
		} else if role.RoleKey == "cc_legion_leader" {
			ccRoleKey = model.RoleTypeLegionLeader
			isCCUser = true
			break
		}
	}

	if !isCCUser {
		return nil
	}

	trace.AddStep(ctx, "Sync CC", "Syncing user %d to CC Member (Role: %s)", user.UserID, ccRoleKey)

	// 3. 检查CC是否已存在
	existingCC, err := s.ccRepo.Get(user.UserID)
	if err == nil {
		// 更新已存在的CC
		existingCC.Name = user.NickName // 优先使用昵称
		if existingCC.Name == "" {
			existingCC.Name = user.UserName
		}
		existingCC.Mobile = user.Phone
		existingCC.RoleType = ccRoleKey
		existingCC.Status = user.Status

		return s.ccRepo.Update(existingCC)
	}

	// 4. 创建新CC
	newCC := &model.CCMember{
		ID:       user.UserID, // 强制使用相同ID
		Name:     user.NickName,
		Mobile:   user.Phone,
		RoleType: ccRoleKey,
		Status:   user.Status,
		// 初始化其他必需字段
		CreateBy: "system", // 或者 user.CreateBy (如果是 *int64 需要处理)
	}
	if newCC.Name == "" {
		newCC.Name = user.UserName
	}
	if newCC.Mobile == "" {
		newCC.Mobile = fmt.Sprintf("Temp%d", user.UserID) // 必须有手机号，临时处理
	}

	return s.ccRepo.Create(newCC)
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	trace.AddStep(ctx, "Start Delete User", "UserID: %d", userID)
	if userID == 1 {
		return errors.New("不允许删除超级管理员")
	}
	trace.AddStep(ctx, "DB Delete", "Deleting user record")
	if err := s.userRepo.Delete(userID); err != nil {
		return err
	}

	// 尝试逻辑删除对应CC
	// 注意：这里可能需要确认需求，删除用户是否也删除CC身份。
	// 暂时假设同步删除
	_ = s.ccRepo.Delete(userID)

	return nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, userID int64, password string) error {
	trace.AddStep(ctx, "Start Reset Password", "UserID: %d", userID)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		trace.AddStep(ctx, "Find User Failed", "User not found")
		return errors.New("用户不存在")
	}
	user.Password = password
	trace.AddStep(ctx, "DB Update Password", "Updating password")
	return s.userRepo.Update(user)
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
