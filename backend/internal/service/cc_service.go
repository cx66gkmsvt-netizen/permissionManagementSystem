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
	// 校验座席号
	if cc.Cno != "" && !s.repo.CheckCnoUnique(cc.Cno, 0) {
		return errors.New("座席号已存在")
	}
	return s.repo.Create(cc)
}

func (s *CCService) Update(cc *model.CCMember) error {
	// 校验手机号
	if !s.repo.CheckMobileUnique(cc.Mobile, cc.ID) {
		return errors.New("手机号已存在")
	}
	// 校验座席号
	if cc.Cno != "" && !s.repo.CheckCnoUnique(cc.Cno, cc.ID) {
		return errors.New("座席号已存在")
	}
	return s.repo.Update(cc)
}

func (s *CCService) Delete(id int64) error {
	return s.repo.Delete(id)
}
