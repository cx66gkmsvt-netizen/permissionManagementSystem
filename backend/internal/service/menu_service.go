package service

import (
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type MenuService struct {
	menuRepo        *repository.MenuRepository
	followUpService *FollowUpService
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo:        repository.NewMenuRepository(),
		followUpService: NewFollowUpService(),
	}
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(menuID int64) (*model.SysMenu, error) {
	return s.menuRepo.FindByID(menuID)
}

// SelectTree 获取菜单树
func (s *MenuService) SelectTree() ([]*model.SysMenu, error) {
	return s.menuRepo.SelectTree()
}

// SelectAll 获取所有菜单
func (s *MenuService) SelectAll() ([]model.SysMenu, error) {
	return s.menuRepo.SelectAll()
}

// Create 创建菜单
func (s *MenuService) Create(req *model.CreateMenuRequest, operatorID int64) error {
	menu := &model.SysMenu{
		ParentID:  req.ParentID,
		MenuName:  req.MenuName,
		MenuType:  req.MenuType,
		Path:      req.Path,
		Component: req.Component,
		Perms:     req.Perms,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Visible:   req.Visible,
		Status:    req.Status,
	}
	if err := s.menuRepo.Create(menu); err != nil {
		return err
	}
	// 记录跟进
	return s.followUpService.Record("sys_menu", menu.MenuID, fmt.Sprintf("创建菜单: %s", menu.MenuName), operatorID, "")
}

// Update 更新菜单
func (s *MenuService) Update(menuID int64, req *model.CreateMenuRequest, operatorID int64) error {
	menu, err := s.menuRepo.FindByID(menuID)
	if err != nil {
		return err
	}

	menu.ParentID = req.ParentID
	menu.MenuName = req.MenuName
	menu.MenuType = req.MenuType
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Perms = req.Perms
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	menu.Status = req.Status

	if err := s.menuRepo.Update(menu); err != nil {
		return err
	}
	// 记录跟进
	return s.followUpService.Record("sys_menu", menuID, "更新菜单信息", operatorID, "")
}

// Delete 删除菜单
func (s *MenuService) Delete(menuID int64, operatorID int64) error {
	if err := s.menuRepo.Delete(menuID); err != nil {
		return err
	}
	// 记录跟进
	return s.followUpService.Record("sys_menu", menuID, "删除菜单", operatorID, "")
}
