package service

import (
	"context"

	"user-center/internal/model"
	"user-center/internal/pkg/trace"
	"user-center/internal/repository"
)

type DeptService struct {
	deptRepo *repository.DeptRepository
}

func NewDeptService() *DeptService {
	return &DeptService{
		deptRepo: repository.NewDeptRepository(),
	}
}

// GetDeptByID 根据ID获取部门
func (s *DeptService) GetDeptByID(deptID int64) (*model.SysDept, error) {
	return s.deptRepo.FindByID(deptID)
}

// SelectTree 获取部门树
func (s *DeptService) SelectTree() ([]*model.SysDept, error) {
	return s.deptRepo.SelectTree()
}

// SelectAll 获取所有部门
func (s *DeptService) SelectAll() ([]model.SysDept, error) {
	return s.deptRepo.SelectAll()
}

// Create 创建部门
func (s *DeptService) Create(ctx context.Context, req *model.CreateDeptRequest) error {
	trace.AddStep(ctx, "Start Create Dept", "DeptName: %s", req.DeptName)
	dept := &model.SysDept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		Sort:     req.Sort,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	trace.AddStep(ctx, "DB Create", "Saving dept")
	return s.deptRepo.Create(dept)
}

// Update 更新部门
func (s *DeptService) Update(ctx context.Context, deptID int64, req *model.CreateDeptRequest) error {
	trace.AddStep(ctx, "Start Update Dept", "DeptID: %d", deptID)
	dept, err := s.deptRepo.FindByID(deptID)
	if err != nil {
		trace.AddStep(ctx, "Find Dept Failed", "Error: %v", err)
		return err
	}

	dept.ParentID = req.ParentID
	dept.DeptName = req.DeptName
	dept.Sort = req.Sort
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email
	dept.Status = req.Status

	trace.AddStep(ctx, "DB Update", "Updating dept")
	return s.deptRepo.Update(dept)
}

// Delete 删除部门
func (s *DeptService) Delete(ctx context.Context, deptID int64) error {
	trace.AddStep(ctx, "Start Delete Dept", "DeptID: %d", deptID)
	return s.deptRepo.Delete(deptID)
}
