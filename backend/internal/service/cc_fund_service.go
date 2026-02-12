package service

import (
	"errors"
	"fmt"
	"user-center/internal/model"
	"user-center/internal/repository"
)

// getRoleTypeName 获取CC角色中文名称
func getRoleTypeName(roleType string) string {
	switch roleType {
	case model.RoleTypeSquadLeader:
		return "CC战队长"
	case model.RoleTypeTeamLeader:
		return "CC团长"
	case model.RoleTypeLegionLeader:
		return "CC军团长"
	default:
		return "CC"
	}
}

type CCFundService struct {
	ccRepo     *repository.CCRepository
	squadRepo  *repository.CCSquadRepository
	teamRepo   *repository.CCTeamRepository
	legionRepo *repository.LegionRepository
	logRepo    *repository.CCLogRepository
}

func NewCCFundService() *CCFundService {
	return &CCFundService{
		ccRepo:     repository.NewCCRepository(),
		squadRepo:  repository.NewCCSquadRepository(),
		teamRepo:   repository.NewCCTeamRepository(),
		legionRepo: repository.NewLegionRepository(),
		logRepo:    repository.NewCCLogRepository(),
	}
}

// ==================== CC资金管理 ====================

// GetCCFund 获取CC资金信息
func (s *CCFundService) GetCCFund(id int64) (map[string]interface{}, error) {
	cc, err := s.ccRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          cc.ID,
		"name":        cc.NickName,
		"balance":     cc.Balance,
		"balanceYuan": float64(cc.Balance) / 100,
	}, nil
}

// EditCCBalance 编辑CC余额
func (s *CCFundService) EditCCBalance(ccID int64, dto *model.FundEditDTO, operatorID int64, operatorName string) error {
	cc, err := s.ccRepo.Get(ccID)
	if err != nil {
		return errors.New("CC不存在")
	}

	var newBalance int64
	if dto.EditType == "set" {
		// 设置为指定值
		if dto.Amount < 0 {
			return errors.New("余额不可以操作至负数")
		}
		newBalance = dto.Amount
	} else {
		// 增减
		newBalance = cc.Balance + dto.Amount
		if newBalance < 0 {
			return errors.New("余额不可以操作至负数")
		}
	}

	oldBalance := cc.Balance
	if err := s.ccRepo.UpdateBalance(ccID, newBalance); err != nil {
		return err
	}

	// 记录日志
	logContent := fmt.Sprintf("修改原因为【%s】从 %.2f元 改为 %.2f元", dto.Reason, float64(oldBalance)/100, float64(newBalance)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeCCBalanceEdit,
		TargetType:    "cc",
		TargetID:      ccID,
		TargetName:    cc.NickName,
		Amount:        newBalance - oldBalance,
		BalanceBefore: oldBalance,
		BalanceAfter:  newBalance,
		Content:       logContent,
		Reason:        dto.Reason,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// GetCCBills 获取CC账单明细
func (s *CCFundService) GetCCBills(ccID int64, billType string) ([]model.CCFundLog, error) {
	return s.logRepo.ListFundLogsByTarget("cc", ccID, billType)
}

// CCTransfer CC转账
func (s *CCFundService) CCTransfer(ccID int64, dto *model.FundTransferDTO, operatorID int64, operatorName string) error {
	cc, err := s.ccRepo.Get(ccID)
	if err != nil {
		return errors.New("CC不存在")
	}

	if cc.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	recipient, err := s.ccRepo.Get(dto.RecipientID)
	if err != nil {
		return errors.New("收账人不存在")
	}

	// 扣除转出方余额
	newCCBalance := cc.Balance - dto.Amount
	if err := s.ccRepo.UpdateBalance(ccID, newCCBalance); err != nil {
		return err
	}

	// 增加收账人余额
	newRecipientBalance := recipient.Balance + dto.Amount
	if err := s.ccRepo.UpdateBalance(dto.RecipientID, newRecipientBalance); err != nil {
		return err
	}

	recipientRole := getRoleTypeName(recipient.RoleType)
	amountYuan := float64(dto.Amount) / 100
	logContent := fmt.Sprintf("%s（CCID:%d），转账%.2f元给%s%s（CCID:%d），个人资金减少%.2f元，对方资金增加%.2f元",
		cc.NickName, ccID, amountYuan, recipientRole, recipient.NickName, recipient.ID, amountYuan, amountYuan)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       "cc_transfer",
		TargetType:    "cc",
		TargetID:      ccID,
		TargetName:    cc.NickName,
		Amount:        -dto.Amount,
		BalanceBefore: cc.Balance,
		BalanceAfter:  newCCBalance,
		RelatedCCID:   &dto.RecipientID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// ==================== 战队资金管理 ====================

// GetSquadFund 获取战队资金信息
func (s *CCFundService) GetSquadFund(id int64) (map[string]interface{}, error) {
	squad, err := s.squadRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          squad.ID,
		"name":        squad.SquadName,
		"balance":     squad.Balance,
		"balanceYuan": float64(squad.Balance) / 100,
	}, nil
}

// EditSquadBalance 编辑战队余额
func (s *CCFundService) EditSquadBalance(squadID int64, dto *model.FundEditDTO, operatorID int64, operatorName string) error {
	squad, err := s.squadRepo.Get(squadID)
	if err != nil {
		return errors.New("战队不存在")
	}

	var newBalance int64
	if dto.EditType == "set" {
		if dto.Amount < 0 {
			return errors.New("余额不可以操作至负数")
		}
		newBalance = dto.Amount
	} else {
		newBalance = squad.Balance + dto.Amount
		if newBalance < 0 {
			return errors.New("余额不可以操作至负数")
		}
	}

	oldBalance := squad.Balance
	if err := s.squadRepo.UpdateBalance(squadID, newBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("修改原因为【%s】从 %.2f元 改为 %.2f元", dto.Reason, float64(oldBalance)/100, float64(newBalance)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeSquadBalanceEdit,
		TargetType:    "squad",
		TargetID:      squadID,
		TargetName:    squad.SquadName,
		Amount:        newBalance - oldBalance,
		BalanceBefore: oldBalance,
		BalanceAfter:  newBalance,
		Content:       logContent,
		Reason:        dto.Reason,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// SquadRecharge 战队充值
func (s *CCFundService) SquadRecharge(squadID int64, dto *model.FundRechargeDTO, operatorID int64, operatorName string) error {
	squad, err := s.squadRepo.Get(squadID)
	if err != nil {
		return errors.New("战队不存在")
	}

	if squad.LeaderID == nil {
		return errors.New("战队没有战队长，无法充值")
	}

	leader, err := s.ccRepo.Get(*squad.LeaderID)
	if err != nil {
		return errors.New("战队长不存在")
	}

	if leader.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	// 扣除战队长余额
	newLeaderBalance := leader.Balance - dto.Amount
	if err := s.ccRepo.UpdateBalance(*squad.LeaderID, newLeaderBalance); err != nil {
		return err
	}

	// 增加战队余额
	newSquadBalance := squad.Balance + dto.Amount
	if err := s.squadRepo.UpdateBalance(squadID, newSquadBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("战队ID%d充值%.2f元，战队长%s（CCID%d）个人资金减少%.2f元，战队资金增加%.2f元",
		squadID, float64(dto.Amount)/100, leader.NickName, leader.ID, float64(dto.Amount)/100, float64(dto.Amount)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeSquadRecharge,
		TargetType:    "squad",
		TargetID:      squadID,
		TargetName:    squad.SquadName,
		Amount:        dto.Amount,
		BalanceBefore: squad.Balance,
		BalanceAfter:  newSquadBalance,
		RelatedCCID:   squad.LeaderID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// SquadTransfer 战队转账
func (s *CCFundService) SquadTransfer(squadID int64, dto *model.FundTransferDTO, operatorID int64, operatorName string) error {
	squad, err := s.squadRepo.Get(squadID)
	if err != nil {
		return errors.New("战队不存在")
	}

	if squad.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	recipient, err := s.ccRepo.Get(dto.RecipientID)
	if err != nil {
		return errors.New("收账人不存在")
	}

	// 扣除战队余额
	newSquadBalance := squad.Balance - dto.Amount
	if err := s.squadRepo.UpdateBalance(squadID, newSquadBalance); err != nil {
		return err
	}

	// 增加收账人余额（按0.9比例）
	transferAmount := int64(float64(dto.Amount) * 0.9)
	newRecipientBalance := recipient.Balance + transferAmount
	if err := s.ccRepo.UpdateBalance(dto.RecipientID, newRecipientBalance); err != nil {
		return err
	}

	recipientRole := getRoleTypeName(recipient.RoleType)
	amountYuan := float64(dto.Amount) / 100
	receivedYuan := float64(transferAmount) / 100
	logContent := fmt.Sprintf("战队名称:%s，战队ID:%d，转账%.2f元给%s%s（CCID:%d），战队资金减少%.2f元，个人资金增加%.2f元",
		squad.SquadName, squadID, amountYuan, recipientRole, recipient.NickName, recipient.ID, amountYuan, receivedYuan)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeSquadTransfer,
		TargetType:    "squad",
		TargetID:      squadID,
		TargetName:    squad.SquadName,
		Amount:        -dto.Amount,
		BalanceBefore: squad.Balance,
		BalanceAfter:  newSquadBalance,
		RelatedCCID:   &dto.RecipientID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// GetSquadBills 获取战队账单明细
func (s *CCFundService) GetSquadBills(squadID int64, billType string) ([]model.CCFundLog, error) {
	return s.logRepo.ListFundLogsByTarget("squad", squadID, billType)
}

// ==================== 团队资金管理 ====================

// GetTeamFund 获取团队资金信息
func (s *CCFundService) GetTeamFund(id int64) (map[string]interface{}, error) {
	team, err := s.teamRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          team.ID,
		"name":        team.TeamName,
		"balance":     team.Balance,
		"balanceYuan": float64(team.Balance) / 100,
	}, nil
}

// EditTeamBalance 编辑团队余额
func (s *CCFundService) EditTeamBalance(teamID int64, dto *model.FundEditDTO, operatorID int64, operatorName string) error {
	team, err := s.teamRepo.Get(teamID)
	if err != nil {
		return errors.New("团队不存在")
	}

	var newBalance int64
	if dto.EditType == "set" {
		if dto.Amount < 0 {
			return errors.New("余额不可以操作至负数")
		}
		newBalance = dto.Amount
	} else {
		newBalance = team.Balance + dto.Amount
		if newBalance < 0 {
			return errors.New("余额不可以操作至负数")
		}
	}

	oldBalance := team.Balance
	if err := s.teamRepo.UpdateBalance(teamID, newBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("修改原因为【%s】从 %.2f元 改为 %.2f元", dto.Reason, float64(oldBalance)/100, float64(newBalance)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeTeamBalanceEdit,
		TargetType:    "team",
		TargetID:      teamID,
		TargetName:    team.TeamName,
		Amount:        newBalance - oldBalance,
		BalanceBefore: oldBalance,
		BalanceAfter:  newBalance,
		Content:       logContent,
		Reason:        dto.Reason,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// TeamRecharge 团队充值
func (s *CCFundService) TeamRecharge(teamID int64, dto *model.FundRechargeDTO, operatorID int64, operatorName string) error {
	team, err := s.teamRepo.Get(teamID)
	if err != nil {
		return errors.New("团队不存在")
	}

	if team.LeaderID == nil {
		return errors.New("团队没有团长，无法充值")
	}

	leader, err := s.ccRepo.Get(*team.LeaderID)
	if err != nil {
		return errors.New("团长不存在")
	}

	if leader.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	newLeaderBalance := leader.Balance - dto.Amount
	if err := s.ccRepo.UpdateBalance(*team.LeaderID, newLeaderBalance); err != nil {
		return err
	}

	newTeamBalance := team.Balance + dto.Amount
	if err := s.teamRepo.UpdateBalance(teamID, newTeamBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("团队ID%d充值%.2f元，团长%s（CCID%d）个人资金减少%.2f元，团队资金增加%.2f元",
		teamID, float64(dto.Amount)/100, leader.NickName, leader.ID, float64(dto.Amount)/100, float64(dto.Amount)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeTeamRecharge,
		TargetType:    "team",
		TargetID:      teamID,
		TargetName:    team.TeamName,
		Amount:        dto.Amount,
		BalanceBefore: team.Balance,
		BalanceAfter:  newTeamBalance,
		RelatedCCID:   team.LeaderID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// TeamTransfer 团队转账
func (s *CCFundService) TeamTransfer(teamID int64, dto *model.FundTransferDTO, operatorID int64, operatorName string) error {
	team, err := s.teamRepo.Get(teamID)
	if err != nil {
		return errors.New("团队不存在")
	}

	if team.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	recipient, err := s.ccRepo.Get(dto.RecipientID)
	if err != nil {
		return errors.New("收账人不存在")
	}

	newTeamBalance := team.Balance - dto.Amount
	if err := s.teamRepo.UpdateBalance(teamID, newTeamBalance); err != nil {
		return err
	}

	transferAmount := int64(float64(dto.Amount) * 0.9)
	newRecipientBalance := recipient.Balance + transferAmount
	if err := s.ccRepo.UpdateBalance(dto.RecipientID, newRecipientBalance); err != nil {
		return err
	}

	recipientRole := getRoleTypeName(recipient.RoleType)
	amountYuan := float64(dto.Amount) / 100
	receivedYuan := float64(transferAmount) / 100
	logContent := fmt.Sprintf("团队名称:%s，团队ID:%d，转账%.2f元给%s%s（CCID:%d），团队资金减少%.2f元，个人资金增加%.2f元",
		team.TeamName, teamID, amountYuan, recipientRole, recipient.NickName, recipient.ID, amountYuan, receivedYuan)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeTeamTransfer,
		TargetType:    "team",
		TargetID:      teamID,
		TargetName:    team.TeamName,
		Amount:        -dto.Amount,
		BalanceBefore: team.Balance,
		BalanceAfter:  newTeamBalance,
		RelatedCCID:   &dto.RecipientID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// GetTeamBills 获取团队账单明细
func (s *CCFundService) GetTeamBills(teamID int64, billType string) ([]model.CCFundLog, error) {
	return s.logRepo.ListFundLogsByTarget("team", teamID, billType)
}

// ==================== 军团资金管理 ====================

// GetLegionFund 获取军团资金信息
func (s *CCFundService) GetLegionFund(id int64) (map[string]interface{}, error) {
	legion, err := s.legionRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          legion.ID,
		"name":        legion.LegionName,
		"balance":     legion.Balance,
		"balanceYuan": float64(legion.Balance) / 100,
	}, nil
}

// EditLegionBalance 编辑军团余额
func (s *CCFundService) EditLegionBalance(legionID int64, dto *model.FundEditDTO, operatorID int64, operatorName string) error {
	legion, err := s.legionRepo.Get(legionID)
	if err != nil {
		return errors.New("军团不存在")
	}

	var newBalance int64
	if dto.EditType == "set" {
		if dto.Amount < 0 {
			return errors.New("余额不可以操作至负数")
		}
		newBalance = dto.Amount
	} else {
		newBalance = legion.Balance + dto.Amount
		if newBalance < 0 {
			return errors.New("余额不可以操作至负数")
		}
	}

	oldBalance := legion.Balance
	if err := s.legionRepo.UpdateBalance(legionID, newBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("修改原因为【%s】从 %.2f元 改为 %.2f元", dto.Reason, float64(oldBalance)/100, float64(newBalance)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeLegionBalanceEdit,
		TargetType:    "legion",
		TargetID:      legionID,
		TargetName:    legion.LegionName,
		Amount:        newBalance - oldBalance,
		BalanceBefore: oldBalance,
		BalanceAfter:  newBalance,
		Content:       logContent,
		Reason:        dto.Reason,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// LegionRecharge 军团充值
func (s *CCFundService) LegionRecharge(legionID int64, dto *model.FundRechargeDTO, operatorID int64, operatorName string) error {
	legion, err := s.legionRepo.Get(legionID)
	if err != nil {
		return errors.New("军团不存在")
	}

	if legion.LeaderID == nil {
		return errors.New("军团没有军团长，无法充值")
	}

	leader, err := s.ccRepo.Get(*legion.LeaderID)
	if err != nil {
		return errors.New("军团长不存在")
	}

	if leader.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	newLeaderBalance := leader.Balance - dto.Amount
	if err := s.ccRepo.UpdateBalance(*legion.LeaderID, newLeaderBalance); err != nil {
		return err
	}

	newLegionBalance := legion.Balance + dto.Amount
	if err := s.legionRepo.UpdateBalance(legionID, newLegionBalance); err != nil {
		return err
	}

	logContent := fmt.Sprintf("军团ID%d充值%.2f元，军团长%s（CCID%d）个人资金减少%.2f元，军团资金增加%.2f元",
		legionID, float64(dto.Amount)/100, leader.NickName, leader.ID, float64(dto.Amount)/100, float64(dto.Amount)/100)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeLegionRecharge,
		TargetType:    "legion",
		TargetID:      legionID,
		TargetName:    legion.LegionName,
		Amount:        dto.Amount,
		BalanceBefore: legion.Balance,
		BalanceAfter:  newLegionBalance,
		RelatedCCID:   legion.LeaderID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// LegionTransfer 军团转账
func (s *CCFundService) LegionTransfer(legionID int64, dto *model.FundTransferDTO, operatorID int64, operatorName string) error {
	legion, err := s.legionRepo.Get(legionID)
	if err != nil {
		return errors.New("军团不存在")
	}

	if legion.Balance < dto.Amount {
		return errors.New("余额不可以操作至负数")
	}

	recipient, err := s.ccRepo.Get(dto.RecipientID)
	if err != nil {
		return errors.New("收账人不存在")
	}

	newLegionBalance := legion.Balance - dto.Amount
	if err := s.legionRepo.UpdateBalance(legionID, newLegionBalance); err != nil {
		return err
	}

	newRecipientBalance := recipient.Balance + dto.Amount
	if err := s.ccRepo.UpdateBalance(dto.RecipientID, newRecipientBalance); err != nil {
		return err
	}

	recipientRole := getRoleTypeName(recipient.RoleType)
	amountYuan := float64(dto.Amount) / 100
	logContent := fmt.Sprintf("军团名称:%s，军团ID:%d，转账%.2f元给%s%s（CCID:%d），军团资金减少%.2f元，个人资金增加%.2f元",
		legion.LegionName, legionID, amountYuan, recipientRole, recipient.NickName, recipient.ID, amountYuan, amountYuan)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeLegionTransfer,
		TargetType:    "legion",
		TargetID:      legionID,
		TargetName:    legion.LegionName,
		Amount:        -dto.Amount,
		BalanceBefore: legion.Balance,
		BalanceAfter:  newLegionBalance,
		RelatedCCID:   &dto.RecipientID,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}

// GetLegionBills 获取军团账单明细
func (s *CCFundService) GetLegionBills(legionID int64, billType string) ([]model.CCFundLog, error) {
	return s.logRepo.ListFundLogsByTarget("legion", legionID, billType)
}

// ==================== 管理员转账 ====================

// AdminTransfer 管理员转账
func (s *CCFundService) AdminTransfer(dto *model.AdminTransferDTO, operatorID int64, operatorName string) error {
	// 获取转出方余额
	var fromBalance int64
	var fromName string
	switch dto.FromType {
	case "legion":
		legion, err := s.legionRepo.Get(dto.FromID)
		if err != nil {
			return errors.New("转出方军团不存在")
		}
		fromBalance = legion.Balance
		fromName = legion.LegionName
	case "team":
		team, err := s.teamRepo.Get(dto.FromID)
		if err != nil {
			return errors.New("转出方团队不存在")
		}
		fromBalance = team.Balance
		fromName = team.TeamName
	case "squad":
		squad, err := s.squadRepo.Get(dto.FromID)
		if err != nil {
			return errors.New("转出方战队不存在")
		}
		fromBalance = squad.Balance
		fromName = squad.SquadName
	case "cc":
		cc, err := s.ccRepo.Get(dto.FromID)
		if err != nil {
			return errors.New("转出方CC不存在")
		}
		fromBalance = cc.Balance
		fromName = cc.NickName
	}

	if fromBalance < dto.Amount {
		return errors.New("余额不足，交易失败")
	}

	// 获取转入方信息
	var toName string
	switch dto.ToType {
	case "legion":
		legion, err := s.legionRepo.Get(dto.ToID)
		if err != nil {
			return errors.New("转入方军团不存在")
		}
		toName = legion.LegionName
	case "team":
		team, err := s.teamRepo.Get(dto.ToID)
		if err != nil {
			return errors.New("转入方团队不存在")
		}
		toName = team.TeamName
	case "squad":
		squad, err := s.squadRepo.Get(dto.ToID)
		if err != nil {
			return errors.New("转入方战队不存在")
		}
		toName = squad.SquadName
	case "cc":
		cc, err := s.ccRepo.Get(dto.ToID)
		if err != nil {
			return errors.New("转入方CC不存在")
		}
		toName = cc.NickName
	}

	// 扣除转出方余额
	switch dto.FromType {
	case "legion":
		s.legionRepo.UpdateBalance(dto.FromID, fromBalance-dto.Amount)
	case "team":
		s.teamRepo.UpdateBalance(dto.FromID, fromBalance-dto.Amount)
	case "squad":
		s.squadRepo.UpdateBalance(dto.FromID, fromBalance-dto.Amount)
	case "cc":
		s.ccRepo.UpdateBalance(dto.FromID, fromBalance-dto.Amount)
	}

	// 增加转入方余额
	switch dto.ToType {
	case "legion":
		legion, _ := s.legionRepo.Get(dto.ToID)
		s.legionRepo.UpdateBalance(dto.ToID, legion.Balance+dto.Amount)
	case "team":
		team, _ := s.teamRepo.Get(dto.ToID)
		s.teamRepo.UpdateBalance(dto.ToID, team.Balance+dto.Amount)
	case "squad":
		squad, _ := s.squadRepo.Get(dto.ToID)
		s.squadRepo.UpdateBalance(dto.ToID, squad.Balance+dto.Amount)
	case "cc":
		cc, _ := s.ccRepo.Get(dto.ToID)
		s.ccRepo.UpdateBalance(dto.ToID, cc.Balance+dto.Amount)
	}

	logContent := fmt.Sprintf("%s ID%d %s 转账%.2f元给 %s ID%d %s",
		dto.FromType, dto.FromID, fromName, float64(dto.Amount)/100, dto.ToType, dto.ToID, toName)
	s.logRepo.CreateFundLog(&model.CCFundLog{
		LogType:       model.FundLogTypeAdminTransfer,
		TargetType:    dto.FromType,
		TargetID:      dto.FromID,
		TargetName:    fromName,
		Amount:        -dto.Amount,
		BalanceBefore: fromBalance,
		BalanceAfter:  fromBalance - dto.Amount,
		Content:       logContent,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
	})

	return nil
}
