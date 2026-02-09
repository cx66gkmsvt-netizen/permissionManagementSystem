package service

import (
	"errors"
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type LegionService struct {
	repo     *repository.LegionRepository
	logRepo  *repository.CCLogRepository
	ccRepo   *repository.CCRepository
	userRepo *repository.UserRepository
}

func NewLegionService() *LegionService {
	return &LegionService{
		repo:     repository.NewLegionRepository(),
		logRepo:  repository.NewCCLogRepository(),
		ccRepo:   repository.NewCCRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// List 军团列表
func (s *LegionService) List(query *model.LegionQuery) (*model.PageResult, error) {
	result, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	// 填充关联信息
	legions := result.Rows.([]model.CCLegion)
	for i := range legions {
		s.fillLegionInfo(&legions[i])
	}
	result.Rows = legions

	return result, nil
}

// ListAll 获取全部军团
func (s *LegionService) ListAll() ([]model.CCLegion, error) {
	return s.repo.ListAll()
}

// Get 获取军团详情
func (s *LegionService) Get(id int64) (*model.CCLegion, error) {
	legion, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	s.fillLegionInfo(legion)
	return legion, nil
}

// Create 创建军团
func (s *LegionService) Create(dto *model.LegionCreateDTO, operatorID int64, operatorName string) error {
	// 校验军团名称唯一
	if !s.repo.CheckNameUnique(dto.LegionName, 0) {
		return errors.New("军团名称不可以重复")
	}

	legion := &model.CCLegion{
		LegionName: dto.LegionName,
		CreateBy:   operatorName,
	}

	if err := s.repo.Create(legion); err != nil {
		return err
	}

	// 记录日志
	logContent := fmt.Sprintf("添加军团 %s（军团ID%d）", legion.LegionName, legion.ID)
	s.logRepo.CreateManageLog(&model.CCManageLog{
		LogType:      model.LogTypeAddLegion,
		TargetType:   "legion",
		TargetID:     legion.ID,
		Content:      logContent,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	})

	return nil
}

// Update 更新军团
func (s *LegionService) Update(id int64, dto *model.LegionUpdateDTO, operatorID int64, operatorName string) error {
	legion, err := s.repo.Get(id)
	if err != nil {
		return errors.New("军团不存在")
	}

	// 校验军团名称唯一
	if !s.repo.CheckNameUnique(dto.LegionName, id) {
		return errors.New("军团名称不可以重复")
	}

	oldLegionName := legion.LegionName
	oldLeaderID := legion.LeaderID

	// 检查是否需要晋升军团长
	leaderChanged := false
	if dto.LeaderID != nil && (oldLeaderID == nil || *dto.LeaderID != *oldLeaderID) {
		leaderChanged = true

		// 验证交易金额
		if dto.TransactionAmount == nil || *dto.TransactionAmount <= 0 {
			return errors.New("晋升军团长需要输入交易金额")
		}

		// 验证新军团长是否存在（从用户管理验证）
		newLeaderUser, err := s.userRepo.FindByID(*dto.LeaderID)
		if err != nil {
			return errors.New("所选军团长不存在")
		}

		// 确保CC记录存在（可能需要创建）
		newLeader, err := s.ccRepo.Get(*dto.LeaderID)
		if err != nil {
			// CC记录不存在，创建一个
			newLeader = &model.CCMember{
				ID:       *dto.LeaderID,
				Name:     newLeaderUser.NickName,
				Mobile:   newLeaderUser.Phone,
				RoleType: model.RoleTypeLegionLeader,
				Status:   "0",
				CreateBy: "system",
			}
			if newLeader.Name == "" {
				newLeader.Name = newLeaderUser.UserName
			}
			if newLeader.Mobile == "" {
				newLeader.Mobile = fmt.Sprintf("Temp%d", *dto.LeaderID)
			}
			if err := s.ccRepo.Create(newLeader); err != nil {
				return errors.New("创建CC记录失败")
			}
			// 重新获取以获取完整信息
			newLeader, _ = s.ccRepo.Get(*dto.LeaderID)
		}

		// 扣除新军团长个人资金
		newBalance := newLeader.Balance - *dto.TransactionAmount
		if err := s.ccRepo.UpdateBalance(*dto.LeaderID, newBalance); err != nil {
			return err
		}

		// 增加军团资金
		if err := s.repo.UpdateBalance(id, legion.Balance+*dto.TransactionAmount); err != nil {
			return err
		}

		// 更新军团长
		if err := s.repo.UpdateLeader(id, dto.LeaderID); err != nil {
			return err
		}

		// 更新CC角色类型
		s.ccRepo.Update(&model.CCMember{ID: *dto.LeaderID, RoleType: model.RoleTypeLegionLeader, LegionID: &id})

		// 记录晋升军团长日志
		fundLogContent := fmt.Sprintf("%s（CCID%d）晋升为%s（军团ID%d）的军团长，扣除个人资金%.2f分，军团资金增加%.2f分",
			newLeader.NickName, newLeader.ID, dto.LegionName, id, float64(*dto.TransactionAmount)/100, float64(*dto.TransactionAmount)/100)
		s.logRepo.CreateFundLog(&model.CCFundLog{
			LogType:         model.FundLogTypePromoteLegionLeader,
			TargetType:      "cc",
			TargetID:        *dto.LeaderID,
			TargetName:      newLeader.NickName,
			Amount:          -*dto.TransactionAmount,
			BalanceBefore:   newLeader.Balance,
			BalanceAfter:    newBalance,
			RelatedLegionID: &id,
			Content:         fundLogContent,
			OperatorID:      operatorID,
			OperatorName:    operatorName,
		})
	}

	// 更新军团名称
	if dto.LegionName != oldLegionName {
		legion.LegionName = dto.LegionName
		legion.UpdateBy = operatorName
		if err := s.repo.Update(legion); err != nil {
			return err
		}

		// 记录修改军团信息日志
		logContent := fmt.Sprintf("修改 %s 的军团名称（军团ID%d），由 %s 改为 %s",
			dto.LegionName, id, oldLegionName, dto.LegionName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyLegionInfo,
			TargetType:   "legion",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	if leaderChanged {
		var oldLeaderName, newLeaderName string
		if oldLeaderID != nil {
			if oldLeader, err := s.ccRepo.Get(*oldLeaderID); err == nil {
				oldLeaderName = oldLeader.NickName
			}
		}
		if dto.LeaderID != nil {
			if newLeader, err := s.ccRepo.Get(*dto.LeaderID); err == nil {
				newLeaderName = newLeader.NickName
			}
		}

		logContent := fmt.Sprintf("修改 %s 的军团长（军团ID%d），由 %s 改为 %s",
			dto.LegionName, id, oldLeaderName, newLeaderName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyLegionInfo,
			TargetType:   "legion",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	return nil
}

// GetLogs 获取军团跟进记录
func (s *LegionService) GetLogs(id int64) ([]model.CCManageLog, error) {
	return s.logRepo.ListManageLogsByLegion(id)
}

// fillLegionInfo 填充军团关联信息
func (s *LegionService) fillLegionInfo(legion *model.CCLegion) {
	if legion.LeaderID != nil {
		if leader, err := s.ccRepo.Get(*legion.LeaderID); err == nil {
			legion.LeaderName = leader.NickName
		}
	}
	// TODO: 计算月度业绩和排名
}
