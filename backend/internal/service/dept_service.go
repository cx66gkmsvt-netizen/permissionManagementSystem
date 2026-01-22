package service

import (
	"user-center/internal/model"
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
func (s *DeptService) Create(req *model.CreateDeptRequest) error {
	dept := &model.SysDept{
		ParentID: req.ParentID,
		DeptName: req.DeptName,
		Sort:     req.Sort,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   req.Status,
	}
	return s.deptRepo.Create(dept)
}

// Update 更新部门
func (s *DeptService) Update(deptID int64, req *model.CreateDeptRequest) error {
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

	return s.deptRepo.Update(dept)
}

// Delete 删除部门
func (s *DeptService) Delete(deptID int64) error {
	return s.deptRepo.Delete(deptID)
}
