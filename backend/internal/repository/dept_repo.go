package repository

import (
	"fmt"
	"strings"
	"user-center/internal/model"

	"gorm.io/gorm"
)

type DeptRepository struct {
	db *gorm.DB
}

func NewDeptRepository() *DeptRepository {
	return &DeptRepository{db: DB}
}

// FindByID 根据ID查找
func (r *DeptRepository) FindByID(deptID int64) (*model.SysDept, error) {
	var dept model.SysDept
	err := r.db.Where("dept_id = ? AND del_flag = '0'", deptID).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// SelectAll 查询所有部门
func (r *DeptRepository) SelectAll() ([]model.SysDept, error) {
	var depts []model.SysDept
	err := r.db.Where("del_flag = '0'").Order("sort ASC").Find(&depts).Error
	return depts, err
}

// SelectTree 获取部门树
func (r *DeptRepository) SelectTree() ([]*model.SysDept, error) {
	depts, err := r.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildDeptTree(depts, 0), nil
}

// buildDeptTree 构建部门树
func buildDeptTree(depts []model.SysDept, parentID int64) []*model.SysDept {
	var tree []*model.SysDept
	for i := range depts {
		if depts[i].ParentID == parentID {
			dept := &depts[i]
			dept.Children = buildDeptTree(depts, dept.DeptID)
			tree = append(tree, dept)
		}
	}
	return tree
}

// Create 创建部门
func (r *DeptRepository) Create(dept *model.SysDept) error {
	// 设置祖级列表
	if dept.ParentID > 0 {
		parent, err := r.FindByID(dept.ParentID)
		if err != nil {
			return err
		}
		dept.Ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.DeptID)
	} else {
		dept.Ancestors = "0"
	}
	return r.db.Create(dept).Error
}

// Update 更新部门
func (r *DeptRepository) Update(dept *model.SysDept) error {
	// 更新祖级列表
	oldDept, err := r.FindByID(dept.DeptID)
	if err != nil {
		return err
	}

	newAncestors := "0"
	if dept.ParentID > 0 {
		parent, err := r.FindByID(dept.ParentID)
		if err != nil {
			return err
		}
		newAncestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.DeptID)
	}
	dept.Ancestors = newAncestors

	// 更新子部门的祖级列表
	if oldDept.Ancestors != newAncestors {
		var children []model.SysDept
		r.db.Where("ancestors LIKE ?", oldDept.Ancestors+",%").Find(&children)
		for _, child := range children {
			child.Ancestors = strings.Replace(child.Ancestors, oldDept.Ancestors, newAncestors, 1)
			r.db.Model(&child).Update("ancestors", child.Ancestors)
		}
	}

	return r.db.Model(dept).Updates(dept).Error
}

// Delete 删除部门(软删除)
func (r *DeptRepository) Delete(deptID int64) error {
	// 检查是否有子部门
	var count int64
	r.db.Model(&model.SysDept{}).Where("parent_id = ? AND del_flag = '0'", deptID).Count(&count)
	if count > 0 {
		return fmt.Errorf("存在子部门，不能删除")
	}
	// 检查是否有用户
	r.db.Model(&model.SysUser{}).Where("dept_id = ? AND del_flag = '0'", deptID).Count(&count)
	if count > 0 {
		return fmt.Errorf("部门下存在用户，不能删除")
	}
	return r.db.Model(&model.SysDept{}).Where("dept_id = ?", deptID).Update("del_flag", "2").Error
}

// GetChildDeptIDs 获取子部门ID列表(包含自己)
func (r *DeptRepository) GetChildDeptIDs(deptID int64) ([]int64, error) {
	var deptIDs []int64
	deptIDs = append(deptIDs, deptID)

	var children []model.SysDept
	err := r.db.Where("FIND_IN_SET(?, ancestors) AND del_flag = '0'", deptID).Find(&children).Error
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		deptIDs = append(deptIDs, child.DeptID)
	}
	return deptIDs, nil
}
