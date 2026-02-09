package service

import (
	"errors"
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

type CCSquadService struct {
	repo       *repository.CCSquadRepository
	teamRepo   *repository.CCTeamRepository
	legionRepo *repository.LegionRepository
	logRepo    *repository.CCLogRepository
	ccRepo     *repository.CCRepository
	userRepo   *repository.UserRepository
}

func NewCCSquadService() *CCSquadService {
	return &CCSquadService{
		repo:       repository.NewCCSquadRepository(),
		teamRepo:   repository.NewCCTeamRepository(),
		legionRepo: repository.NewLegionRepository(),
		logRepo:    repository.NewCCLogRepository(),
		ccRepo:     repository.NewCCRepository(),
		userRepo:   repository.NewUserRepository(),
	}
}

// List 战队列表
func (s *CCSquadService) List(query *model.SquadQuery) (*model.PageResult, error) {
	result, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	squads := result.Rows.([]model.CCSquad)
	for i := range squads {
		s.fillSquadInfo(&squads[i])
	}
	result.Rows = squads

	return result, nil
}

// ListAll 获取全部战队
func (s *CCSquadService) ListAll() ([]model.CCSquad, error) {
	return s.repo.ListAll()
}

// ListByTeamID 根据团队ID获取战队列表
func (s *CCSquadService) ListByTeamID(teamID int64) ([]model.CCSquad, error) {
	return s.repo.ListByTeamID(teamID)
}

// Get 获取战队详情
func (s *CCSquadService) Get(id int64) (*model.CCSquad, error) {
	squad, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	s.fillSquadInfo(squad)
	return squad, nil
}

// Create 创建战队
func (s *CCSquadService) Create(dto *model.SquadCreateDTO, operatorID int64, operatorName string) error {
	// 校验战队名称唯一
	if !s.repo.CheckNameUnique(dto.SquadName, 0) {
		return errors.New("战队名称不可重复")
	}

	// 获取所属团队
	team, err := s.teamRepo.Get(dto.TeamID)
	if err != nil {
		return errors.New("所选团队不存在")
	}

	// 扣除团队资金
	if dto.TransactionAmount > 0 {
		if team.Balance < dto.TransactionAmount {
			return errors.New("余额不可以操作至负数")
		}

		newBalance := team.Balance - dto.TransactionAmount
		if err := s.teamRepo.UpdateBalance(dto.TeamID, newBalance); err != nil {
			return err
		}
	}

	squad := &model.CCSquad{
		SquadName: dto.SquadName,
		TeamID:    dto.TeamID,
		CreateBy:  operatorName,
	}

	if err := s.repo.Create(squad); err != nil {
		return err
	}

	// 记录管理日志
	logContent := fmt.Sprintf("添加战队 %s（战队ID%d）所属团队为 %s（团队ID%d）",
		squad.SquadName, squad.ID, team.TeamName, team.ID)
	s.logRepo.CreateManageLog(&model.CCManageLog{
		LogType:      model.LogTypeAddSquad,
		TargetType:   "squad",
		TargetID:     squad.ID,
		Content:      logContent,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	})

	// 记录资金日志
	if dto.TransactionAmount > 0 {
		fundLogContent := fmt.Sprintf("%s（团队ID%d）添加战队%s（战队ID%d）减少团队资金%.2f分",
			team.TeamName, team.ID, squad.SquadName, squad.ID, float64(dto.TransactionAmount)/100)
		s.logRepo.CreateFundLog(&model.CCFundLog{
			LogType:        model.FundLogTypeAddSquad,
			TargetType:     "team",
			TargetID:       dto.TeamID,
			TargetName:     team.TeamName,
			Amount:         -dto.TransactionAmount,
			BalanceBefore:  team.Balance,
			BalanceAfter:   team.Balance - dto.TransactionAmount,
			RelatedSquadID: &squad.ID,
			Content:        fundLogContent,
			OperatorID:     operatorID,
			OperatorName:   operatorName,
		})
	}

	return nil
}

// Update 更新战队
func (s *CCSquadService) Update(id int64, dto *model.SquadUpdateDTO, operatorID int64, operatorName string) error {
	squad, err := s.repo.Get(id)
	if err != nil {
		return errors.New("战队不存在")
	}

	// 校验战队名称唯一
	if !s.repo.CheckNameUnique(dto.SquadName, id) {
		return errors.New("战队名称不可重复")
	}

	oldSquadName := squad.SquadName
	oldLeaderID := squad.LeaderID
	oldTeamID := squad.TeamID

	// 检查战队长是否变更
	leaderChanged := false
	if dto.LeaderID != nil && (oldLeaderID == nil || *dto.LeaderID != *oldLeaderID) {
		leaderChanged = true

		// 验证新战队长是否已担任其他战队战队长
		if !s.repo.CheckLeaderUnique(*dto.LeaderID, id) {
			return errors.New("该用户已担任其他战队的战队长，不可选择")
		}

		// 验证交易金额
		if dto.TransactionAmount == nil || *dto.TransactionAmount <= 0 {
			return errors.New("晋升战队长需要输入交易金额")
		}

		// 验证新战队长是否存在（从用户管理验证）
		newLeaderUser, err := s.userRepo.FindByID(*dto.LeaderID)
		if err != nil {
			return errors.New("所选战队长不存在")
		}

		// 确保CC记录存在（可能需要创建）
		newLeader, err := s.ccRepo.Get(*dto.LeaderID)
		if err != nil {
			// CC记录不存在，创建一个
			newLeader = &model.CCMember{
				ID:       *dto.LeaderID,
				Name:     newLeaderUser.NickName,
				Mobile:   newLeaderUser.Phone,
				RoleType: model.RoleTypeSquadLeader,
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

		// 扣除新战队长个人资金
		newBalance := newLeader.Balance - *dto.TransactionAmount
		if err := s.ccRepo.UpdateBalance(*dto.LeaderID, newBalance); err != nil {
			return err
		}

		// 增加战队资金
		if err := s.repo.UpdateBalance(id, squad.Balance+*dto.TransactionAmount); err != nil {
			return err
		}

		// 更新战队长
		if err := s.repo.UpdateLeader(id, dto.LeaderID); err != nil {
			return err
		}

		// 更新CC角色类型
		s.ccRepo.Update(&model.CCMember{ID: *dto.LeaderID, RoleType: model.RoleTypeSquadLeader, SquadID: &id})

		// 记录晋升战队长日志
		fundLogContent := fmt.Sprintf("%s（CCID%d）晋升为%s（战队ID%d）的战队长，扣除个人资金%.2f分，战队资金增加%.2f分",
			newLeader.NickName, newLeader.ID, dto.SquadName, id, float64(*dto.TransactionAmount)/100, float64(*dto.TransactionAmount)/100)
		s.logRepo.CreateFundLog(&model.CCFundLog{
			LogType:        model.FundLogTypePromoteSquadLeader,
			TargetType:     "cc",
			TargetID:       *dto.LeaderID,
			TargetName:     newLeader.NickName,
			Amount:         -*dto.TransactionAmount,
			BalanceBefore:  newLeader.Balance,
			BalanceAfter:   newBalance,
			RelatedSquadID: &id,
			Content:        fundLogContent,
			OperatorID:     operatorID,
			OperatorName:   operatorName,
		})
	}

	// 检查所属团队是否变更
	teamChanged := false
	if dto.TeamID != oldTeamID {
		teamChanged = true
		squad.TeamID = dto.TeamID
	}

	// 更新基本信息
	squad.SquadName = dto.SquadName
	squad.UpdateBy = operatorName
	if err := s.repo.Update(squad); err != nil {
		return err
	}

	// 记录名称变更日志
	if dto.SquadName != oldSquadName {
		logContent := fmt.Sprintf("修改 战队名称（战队ID%d），由 %s 改为 %s",
			id, oldSquadName, dto.SquadName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifySquadInfo,
			TargetType:   "squad",
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

		logContent := fmt.Sprintf("修改 %s 的战队长（战队ID%d），由 %s 改为 %s",
			dto.SquadName, id, oldLeaderName, newLeaderName)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifySquadStruct,
			TargetType:   "squad",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	if teamChanged {
		var oldTeamName, newTeamName string
		if oldTeam, err := s.teamRepo.Get(oldTeamID); err == nil {
			oldTeamName = oldTeam.TeamName
		}
		if newTeam, err := s.teamRepo.Get(dto.TeamID); err == nil {
			newTeamName = newTeam.TeamName
		}

		logContent := fmt.Sprintf("修改 %s 的所属团队，由 %s（团队ID%d）改为 %s（团队ID%d）",
			dto.SquadName, oldTeamName, oldTeamID, newTeamName, dto.TeamID)
		s.logRepo.CreateManageLog(&model.CCManageLog{
			LogType:      model.LogTypeModifySquadStruct,
			TargetType:   "squad",
			TargetID:     id,
			Content:      logContent,
			OperatorID:   operatorID,
			OperatorName: operatorName,
		})
	}

	return nil
}

// GetLogs 获取战队修改记录
func (s *CCSquadService) GetLogs(id int64) ([]model.CCManageLog, error) {
	return s.logRepo.ListManageLogsBySquad(id)
}

// fillSquadInfo 填充战队关联信息
func (s *CCSquadService) fillSquadInfo(squad *model.CCSquad) {
	if squad.LeaderID != nil {
		if leader, err := s.userRepo.FindByID(*squad.LeaderID); err == nil {
			squad.LeaderName = leader.NickName
			if squad.LeaderName == "" {
				squad.LeaderName = leader.UserName
			}
		}
	}

	if team, err := s.teamRepo.Get(squad.TeamID); err == nil {
		squad.TeamName = team.TeamName
		if team.LeaderID != nil {
			if teamLeader, err := s.userRepo.FindByID(*team.LeaderID); err == nil {
				squad.TeamLeaderName = teamLeader.NickName
				if squad.TeamLeaderName == "" {
					squad.TeamLeaderName = teamLeader.UserName
				}
			}
		}

		if team.LegionID != nil {
			squad.LegionID = team.LegionID
			if legion, err := s.legionRepo.Get(*team.LegionID); err == nil {
				squad.LegionName = legion.LegionName
				if legion.LeaderID != nil {
					if legionLeader, err := s.userRepo.FindByID(*legion.LeaderID); err == nil {
						squad.LegionLeaderName = legionLeader.NickName
						if squad.LegionLeaderName == "" {
							squad.LegionLeaderName = legionLeader.UserName
						}
					}
				}
			}
		}
	}
}
