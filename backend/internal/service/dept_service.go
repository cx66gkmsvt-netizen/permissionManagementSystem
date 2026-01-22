package service

import (
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type DeptService struct {
	deptRepo        *repository.DeptRepository
	followUpService *FollowUpService
}

func NewDeptService() *DeptService {
	return &DeptService{
		deptRepo:        repository.NewDeptRepository(),
		followUpService: NewFollowUpService(),
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
func (s *DeptService) Create(req *model.CreateDeptRequest, operatorID int64) error {
	dept := &model.SysDept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		Sort:     req.Sort,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	if err := s.deptRepo.Create(dept); err != nil {
		return err
	}

	// 记录跟进
	return s.followUpService.Record("sys_dept", dept.DeptID, fmt.Sprintf("创建部门: %s", dept.DeptName), operatorID, "")
}

// Update 更新部门
func (s *DeptService) Update(deptID int64, req *model.CreateDeptRequest, operatorID int64) error {
	dept, err := s.deptRepo.FindByID(deptID)
	if err != nil {
		return err
	}

	dept.ParentID = req.ParentID
	dept.DeptName = req.DeptName
	dept.Sort = req.Sort
	dept.Leader = req.Leader
	dept.Phone = req.Phone
	dept.Email = req.Email
	dept.Status = req.Status

	if err := s.deptRepo.Update(dept); err != nil {
		return err
	}

	// 记录跟进
	return s.followUpService.Record("sys_dept", deptID, "更新部门信息", operatorID, "")
}

// Delete 删除部门
func (s *DeptService) Delete(deptID int64, operatorID int64) error {
	if err := s.deptRepo.Delete(deptID); err != nil {
		return err
	}
	// 记录跟进
	return s.followUpService.Record("sys_dept", deptID, "删除部门", operatorID, "")
}
