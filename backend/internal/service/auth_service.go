package service

import (
	"errors"

	"user-center/internal/config"
	"user-center/internal/model"
	"user-center/internal/pkg"
	"user-center/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	menuRepo *repository.MenuRepository
	cfg      *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
		menuRepo: repository.NewMenuRepository(),
		cfg:      cfg,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *model.LoginRequest, ip string) (*model.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUserName(req.UserName)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查状态
	if user.Status == "1" {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !pkg.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新登录信息
	s.userRepo.UpdateLoginInfo(user.UserID, ip)

	// 生成Token
	token, expireAt, err := pkg.GenerateToken(user.UserID, user.UserName, s.cfg.JWT.ExpireTime)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	return &model.LoginResponse{
		Token:    token,
		ExpireAt: expireAt,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID int64) (*model.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 获取角色Key列表
	var roleKeys []string
	for _, role := range user.Roles {
		roleKeys = append(roleKeys, role.RoleKey)
	}

	// 获取权限列表
	var perms []string
	if userID == 1 {
		perms = []string{"*:*:*"} // 超管拥有所有权限
	} else {
		var roleIDs []int64
		for _, role := range user.Roles {
			roleIDs = append(roleIDs, role.RoleID)
		}
		perms, _ = s.menuRepo.SelectPermsByRoleIDs(roleIDs)
	}

	return &model.UserInfo{
		User:        user,
		Roles:       roleKeys,
		Permissions: perms,
	}, nil
}

// GetRoutes 获取路由菜单
func (s *AuthService) GetRoutes(userID int64) ([]*model.SysMenu, error) {
	var menus []model.SysMenu
	var err error

	if userID == 1 {
		// 超管获取所有菜单
		menus, err = s.menuRepo.SelectAll()
	} else {
		// 根据角色获取菜单
		roles, _ := s.userRepo.GetUserRoles(userID)
		var roleIDs []int64
		for _, role := range roles {
			roleIDs = append(roleIDs, role.RoleID)
		}
		menus, err = s.menuRepo.SelectMenusByRoleIDs(roleIDs)
	}

	if err != nil {
		return nil, err
	}

	// 过滤出目录和菜单(不含按钮)
	var filtered []model.SysMenu
	for _, menu := range menus {
		if menu.MenuType != "F" {
			filtered = append(filtered, menu)
		}
	}

	return buildMenuTree(filtered, 0), nil
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
