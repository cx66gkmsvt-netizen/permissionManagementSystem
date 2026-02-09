package service

import (
	"errors"
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type CCTeamService struct {
	repo       *repository.CCTeamRepository
	legionRepo *repository.LegionRepository
	logRepo    *repository.CCLogRepository
	ccRepo     *repository.CCRepository
	userRepo   *repository.UserRepository
}

func NewCCTeamService() *CCTeamService {
	return &CCTeamService{
		repo:       repository.NewCCTeamRepository(),
		legionRepo: repository.NewLegionRepository(),
		logRepo:    repository.NewCCLogRepository(),
		ccRepo:     repository.NewCCRepository(),
		userRepo:   repository.NewUserRepository(),
	}
}

// List 团队列表
func (s *CCTeamService) List(query *model.TeamQuery) (*model.PageResult, error) {
	result, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	teams := result.Rows.([]model.CCTeam)
	for i := range teams {
		s.fillTeamInfo(&teams[i])
	}
	result.Rows = teams

	return result, nil
}

// ListAll 获取全部团队
func (s *CCTeamService) ListAll() ([]model.CCTeam, error) {
	return s.repo.ListAll()
}

// Get 获取团队详情
func (s *CCTeamService) Get(id int64) (*model.CCTeam, error) {
	team, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	s.fillTeamInfo(team)
	return team, nil
}

// Create 创建团队
func (s *CCTeamService) Create(dto *model.TeamCreateDTO, operatorID int64, operatorName string) error {
	team := &model.CCTeam{
		TeamName:     dto.TeamName,
		BusinessType: dto.BusinessType,
		LegionID:     dto.LegionID,
		CreateBy:     operatorName,
	}

	// 如果有所属军团，需要扣除军团资金
	if dto.LegionID != nil && dto.TransactionAmount > 0 {
		legion, err := s.legionRepo.Get(*dto.LegionID)
		if err != nil {
			return errors.New("所选军团不存在")
		}

		if legion.Balance < dto.TransactionAmount {
			return errors.New("余额不可以操作至负数")
		}

		// 扣除军团资金
		newBalance := legion.Balance - dto.TransactionAmount
		if err := s.legionRepo.UpdateBalance(*dto.LegionID, newBalance); err != nil {
			return err
		}
	}

	if err := s.repo.Create(team); err != nil {
		return err
	}

	// 记录日志
	var legionInfo string
	if dto.LegionID != nil {
		if legion, err := s.legionRepo.Get(*dto.LegionID); err == nil {
			legionInfo = fmt.Sprintf("所属军团为 %s（军团ID%d）", legion.LegionName, legion.ID)
		}
	} else {
		legionInfo = "无所属军团"
	}
	logContent := fmt.Sprintf("添加团队 %s（团队ID%d），业务类型为 %s，%s",
		team.TeamName, team.ID, dto.BusinessType, legionInfo)
	s.logRepo.CreateManageLog(&model.CCManageLog{
		LogType:      model.LogTypeAddTeam,
		TargetType:   "team",
		TargetID:     team.ID,
		Content:      logContent,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	})

	// 记录资金日志
	if dto.LegionID != nil && dto.TransactionAmount > 0 {
		legion, _ := s.legionRepo.Get(*dto.LegionID)
		fundLogContent := fmt.Sprintf("%s（军团ID%d）添加团队%s（团队ID%d）减少军团资金%.2f分",
			legion.LegionName, legion.ID, team.TeamName, team.ID, float64(dto.TransactionAmount)/100)
		s.logRepo.CreateFundLog(&model.CCFundLog{
			LogType:       model.FundLogTypeAddTeam,
			TargetType:    "legion",
			TargetID:      *dto.LegionID,
			TargetName:    legion.LegionName,
			Amount:        -dto.TransactionAmount,
			BalanceBefore: legion.Balance + dto.TransactionAmount,
			BalanceAfter:  legion.Balance,
			RelatedTeamID: &team.ID,
			Content:       fundLogContent,
			OperatorID:    operatorID,
			OperatorName:  operatorName,
		})
	}

	return nil
}

// Update 更新团队
func (s *CCTeamService) Update(id int64, dto *model.TeamUpdateDTO, operatorID int64, operatorName string) error {
	team, err := s.repo.Get(id)
	if err != nil {
		return errors.New("团队不存在")
	}

	oldTeamName := team.TeamName
	oldBusinessType := team.BusinessType
	oldLeaderID := team.LeaderID
	oldLegionID := team.LegionID

	// 检查团长是否变更
	leaderChanged := false
	if dto.LeaderID != nil && (oldLeaderID == nil || *dto.LeaderID != *oldLeaderID) {
		leaderChanged = true

		// 验证交易金额
		if dto.TransactionAmount == nil || *dto.TransactionAmount <= 0 {
			return errors.New("晋升团长需要输入交易金额")
		}

		// 验证新团长是否存在（从用户管理验证）
		newLeaderUser, err := s.userRepo.FindByID(*dto.LeaderID)
		if err != nil {
			return errors.New("所选团长不存在")
		}

		// 确保CC记录存在（可能需要创建）
		newLeader, err := s.ccRepo.Get(*dto.LeaderID)
		if err != nil {
			// CC记录不存在，创建一个
			newLeader = &model.CCMember{
				ID:       *dto.LeaderID,
				Name:     newLeaderUser.NickName,
				Mobile:   newLeaderUser.Phone,
				RoleType: model.RoleTypeTeamLeader,
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
			newLeader, _ = s.ccRepo.Get(*dto.LeaderID)
		}

		// 扣除新团长个人资金
		newBalance := newLeader.Balance - *dto.TransactionAmount
		if err := s.ccRepo.UpdateBalance(*dto.LeaderID, newBalance); err != nil {
			return err
		}

		// 增加团队资金
		if err := s.repo.UpdateBalance(id, team.Balance+*dto.TransactionAmount); err != nil {
			return err
		}

		// 更新团长
		if err := s.repo.UpdateLeader(id, dto.LeaderID); err != nil {
			return err
		}

		// 更新CC角色类型
		s.ccRepo.Update(&model.CCMember{ID: *dto.LeaderID, RoleType: model.RoleTypeTeamLeader, TeamID: &id})

		// 记录晋升团长日志
		fundLogContent := fmt.Sprintf("%s（CCID%d）晋升为%s（团队ID%d）的团长，扣除个人资金%.2f分，团队资金增加%.2f分",
			newLeader.NickName, newLeader.ID, dto.TeamName, id, float64(*dto.TransactionAmount)/100, float64(*dto.TransactionAmount)/100)
		s.logRepo.CreateFundLog(&model.CCFundLog{
			LogType:       model.FundLogTypePromoteTeamLeader,
			TargetType:    "cc",
			TargetID:      *dto.LeaderID,
			TargetName:    newLeader.NickName,
			Amount:        -*dto.TransactionAmount,
			BalanceBefore: newLeader.Balance,
			BalanceAfter:  newBalance,
			RelatedTeamID: &id,
			Content:       fundLogContent,
			OperatorID:    operatorID,
			OperatorName:  operatorName,
		})
	}

	// 检查所属军团是否变更
	legionChanged := false
	if dto.LegionID != nil && (oldLegionID == nil || *dto.LegionID != *oldLegionID) {
		legionChanged = true
		team.LegionID = dto.LegionID
	} else if dto.LegionID == nil && oldLegionID != nil {
		legionChanged = true
		team.LegionID = nil
	}

	// 更新基本信息
	team.TeamName = dto.TeamName
	team.BusinessType = dto.BusinessType
	team.UpdateBy = operatorName
	if err := s.repo.Update(team); err != nil {
		return err
	}

	// 记录名称/业务类型变更日志
	if dto.TeamName != oldTeamName {
		logContent := fmt.Sprintf("修改 %s 的团队名称（团队ID%d），由 %s 改为 %s",
			dto.TeamName, id, oldTeamName, dto.TeamName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyTeamInfo,
			TargetType:   "team",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	if dto.BusinessType != oldBusinessType {
		logContent := fmt.Sprintf("修改 %s 的业务类型（团队ID%d），由 %s 改为 %s",
			dto.TeamName, id, oldBusinessType, dto.BusinessType)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyTeamInfo,
			TargetType:   "team",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	if leaderChanged {
		var oldLeaderName, newLeaderName string
		if oldLeaderID != nil {
			if oldLeaderUser, err := s.userRepo.FindByID(*oldLeaderID); err == nil {
				oldLeaderName = oldLeaderUser.NickName
				if oldLeaderName == "" {
					oldLeaderName = oldLeaderUser.UserName
				}
			}
		}
		if dto.LeaderID != nil {
			if newLeaderUser, err := s.userRepo.FindByID(*dto.LeaderID); err == nil {
				newLeaderName = newLeaderUser.NickName
				if newLeaderName == "" {
					newLeaderName = newLeaderUser.UserName
				}
			}
		}

		logContent := fmt.Sprintf("修改 %s 的团长（团队ID%d），由 %s 改为 %s",
			dto.TeamName, id, oldLeaderName, newLeaderName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyTeamInfo,
			TargetType:   "team",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	if legionChanged {
		var oldLegionName, newLegionName string
		if oldLegionID != nil {
			if oldLegion, err := s.legionRepo.Get(*oldLegionID); err == nil {
				oldLegionName = oldLegion.LegionName
			}
		}
		if dto.LegionID != nil {
			if newLegion, err := s.legionRepo.Get(*dto.LegionID); err == nil {
				newLegionName = newLegion.LegionName
			}
		}

		logContent := fmt.Sprintf("修改 %s 的所属军团，由 %s 改为 %s",
			dto.TeamName, oldLegionName, newLegionName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifyTeamStruct,
			TargetType:   "team",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	return nil
}

// GetLogs 获取团队修改记录
func (s *CCTeamService) GetLogs(id int64) ([]model.CCManageLog, error) {
	return s.logRepo.ListManageLogsByTeam(id)
}

// fillTeamInfo 填充团队关联信息
func (s *CCTeamService) fillTeamInfo(team *model.CCTeam) {
	if team.LeaderID != nil {
		if leader, err := s.userRepo.FindByID(*team.LeaderID); err == nil {
			team.LeaderName = leader.NickName
			if team.LeaderName == "" {
				team.LeaderName = leader.UserName
			}
		}
	}
	if team.LegionID != nil {
		if legion, err := s.legionRepo.Get(*team.LegionID); err == nil {
			team.LegionName = legion.LegionName
			if legion.LeaderID != nil {
				if legionLeader, err := s.userRepo.FindByID(*legion.LeaderID); err == nil {
					team.LegionLeaderName = legionLeader.NickName
					if team.LegionLeaderName == "" {
						team.LegionLeaderName = legionLeader.UserName
					}
				}
			}
		}
	}
}
