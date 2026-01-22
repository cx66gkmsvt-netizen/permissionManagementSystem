package service

import (
	"fmt"
	"time"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type FollowUpService struct {
	followUpRepo *repository.FollowUpRepository
	userRepo     *repository.UserRepository
}

func NewFollowUpService() *FollowUpService {
	return &FollowUpService{
		followUpRepo: repository.NewFollowUpRepository(),
		userRepo:     repository.NewUserRepository(), // To get operator name
	}
}

// Record 记录跟进信息
func (s *FollowUpService) Record(targetType string, targetID int64, content string, operatorID int64, remark string) error {
	// 获取操作人名称
	var operName string
	user, err := s.userRepo.FindByID(operatorID)
	if err == nil && user != nil {
		operName = user.UserName // 或者 NickName
	} else {
		// 如果找不到用户，或者是系统操作，或者是 ID=0 (e.g. initial seed)
		if operatorID == 1 { // Assuming 1 is admin/system usually
			operName = "admin"
		} else {
			operName = fmt.Sprintf("Unknown(%d)", operatorID)
		}
	}

	record := &model.SysFollowUp{
		TargetType:   targetType,
		TargetID:     targetID,
		Content:      content,
		OperUserID:   &operatorID,
		OperUserName: operName,
		OperTime:     time.Now(),
		Remark:       remark,
	}

	return s.followUpRepo.Create(record)
}
