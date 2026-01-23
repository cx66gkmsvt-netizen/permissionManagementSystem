package service

import (
	"context"

	"user-center/internal/model"
	"user-center/internal/pkg/trace"
	"user-center/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo: repository.NewMenuRepository(),
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
func (s *MenuService) Create(ctx context.Context, req *model.CreateMenuRequest) error {
	trace.AddStep(ctx, "Start Create Menu", "MenuName: %s", req.MenuName)
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
	trace.AddStep(ctx, "DB Create", "Saving menu")
	return s.menuRepo.Create(menu)
}

// Update 更新菜单
func (s *MenuService) Update(ctx context.Context, menuID int64, req *model.CreateMenuRequest) error {
	trace.AddStep(ctx, "Start Update Menu", "MenuID: %d", menuID)
	menu, err := s.menuRepo.FindByID(menuID)
	if err != nil {
		trace.AddStep(ctx, "Find Menu Failed", "Error: %v", err)
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

	trace.AddStep(ctx, "DB Update", "Updating menu")
	return s.menuRepo.Update(menu)
}

// Delete 删除菜单
func (s *MenuService) Delete(ctx context.Context, menuID int64) error {
	trace.AddStep(ctx, "Start Delete Menu", "MenuID: %d", menuID)
	return s.menuRepo.Delete(menuID)
}
