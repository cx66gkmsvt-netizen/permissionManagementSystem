package service

import (
	"errors"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type CCService struct {
	repo *repository.CCRepository
}

func NewCCService() *CCService {
	return &CCService{
		repo: repository.NewCCRepository(),
	}
}

func (s *CCService) List(query *model.CCQuery) (*model.PageResult, error) {
	return s.repo.List(query)
}

func (s *CCService) Get(id int64) (*model.CCMember, error) {
	return s.repo.Get(id)
}

func (s *CCService) Create(cc *model.CCMember) error {
	// 校验手机号
	if !s.repo.CheckMobileUnique(cc.Mobile, 0) {
		return errors.New("手机号已存在")
	}
	// 校验容联座席号
	if cc.RonglianSeat != "" && !s.repo.CheckRonglianSeatUnique(cc.RonglianSeat, 0) {
		return errors.New("容联座席号已存在")
	}
	// 设置默认角色类型
	if cc.RoleType == "" {
		cc.RoleType = model.RoleTypeCC
	}
	return s.repo.Create(cc)
}

func (s *CCService) Update(cc *model.CCMember) error {
	// 校验手机号
	if cc.Mobile != "" && !s.repo.CheckMobileUnique(cc.Mobile, cc.ID) {
		return errors.New("手机号已存在")
	}
	// 校验容联座席号
	if cc.RonglianSeat != "" && !s.repo.CheckRonglianSeatUnique(cc.RonglianSeat, cc.ID) {
		return errors.New("容联座席号已存在")
	}
	// 校验云客账号1和云客账号2不能相同
	if cc.CloudAccount1 != "" && cc.CloudAccount2 != "" && cc.CloudAccount1 == cc.CloudAccount2 {
		return errors.New("CC云客账号1与云客账号2重复，请修改后再填写")
	}
	return s.repo.Update(cc)
}

func (s *CCService) Delete(id int64) error {
	return s.repo.Delete(id)
}
