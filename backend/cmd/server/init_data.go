package main

import (
	"log"
	"time"
	"user-center/internal/model"
	"user-center/internal/repository"
)

// initCCData 初始化CC相关数据
func initCCData() {
	db := repository.GetDB()

	// 检查是否已有军团数据
	var count int64
	db.Model(&model.CCLegion{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("Initializing CC data...")

	// 1. 创建CC成员（如果不复用sys_user，则单独创建）
	// 这里为了简单，假设没有和SysUser强关联，或者已经同步。
	// 但之前的UserHandler同步逻辑是：创建SysUser -> 创建CCMember。
	// 为了测试方便，我们创建几个独立的CCMember作为各级Leader。

	// 军团长
	legionLeader := &model.CCMember{
		Name:             "Legion Leader",
		NickName:         "军团长A",
		Mobile:           "13800138001",
		RoleType:         model.RoleTypeLegionLeader,
		Balance:          1000000, // 10000元
		Status:           "0",
		AttendanceStatus: model.AttendanceStatusOnDuty,
	}
	db.Create(legionLeader)

	// 团长
	teamLeader := &model.CCMember{
		Name:             "Team Leader",
		NickName:         "团长B",
		Mobile:           "13800138002",
		RoleType:         model.RoleTypeTeamLeader,
		Balance:          500000, // 5000元
		Status:           "0",
		AttendanceStatus: model.AttendanceStatusOnDuty,
	}
	db.Create(teamLeader)

	// 战队长
	squadLeader := &model.CCMember{
		Name:             "Squad Leader",
		NickName:         "战队长C",
		Mobile:           "13800138003",
		RoleType:         model.RoleTypeSquadLeader,
		Balance:          200000, // 2000元
		Status:           "0",
		AttendanceStatus: model.AttendanceStatusOnDuty,
	}
	db.Create(squadLeader)

	// 普通CC
	cc := &model.CCMember{
		Name:             "CC Member",
		NickName:         "CC小张",
		Mobile:           "13800138004",
		RoleType:         model.RoleTypeCC,
		Balance:          10000, // 100元
		Status:           "0",
		AttendanceStatus: model.AttendanceStatusOnDuty,
	}
	db.Create(cc)

	// 2. 创建军团
	legion := &model.CCLegion{
		LegionName: "第一军团",
		LeaderID:   &legionLeader.ID,
		Balance:    10000000, // 10万元
		Status:     "0",
		CreateBy:   "admin",
	}
	db.Create(legion)

	// 更新军团长的LegionID
	db.Model(legionLeader).Update("legion_id", legion.ID)

	// 3. 创建团队
	team := &model.CCTeam{
		TeamName:     "先锋团队",
		BusinessType: "业务A",
		LeaderID:     &teamLeader.ID,
		LegionID:     &legion.ID,
		Balance:      5000000, // 5万元
		Status:       "0",
		CreateBy:     "admin",
	}
	db.Create(team)

	// 更新团长的TeamID和LegionID
	db.Model(teamLeader).Updates(map[string]interface{}{
		"team_id":   team.ID,
		"legion_id": legion.ID,
	})

	// 4. 创建战队
	squad := &model.CCSquad{
		SquadName: "突击战队",
		LeaderID:  &squadLeader.ID,
		TeamID:    team.ID,
		Balance:   1000000, // 1万元
		Status:    "0",
		CreateBy:  "admin",
	}
	db.Create(squad)

	// 更新战队长的SquadID, TeamID, LegionID
	db.Model(squadLeader).Updates(map[string]interface{}{
		"squad_id":  squad.ID,
		"team_id":   team.ID,
		"legion_id": legion.ID,
	})

	// 更新普通CC的SquadID, TeamID, LegionID
	db.Model(cc).Updates(map[string]interface{}{
		"squad_id":  squad.ID,
		"team_id":   team.ID,
		"legion_id": legion.ID,
	})

	// 5. 创建一些日志记录
	db.Create(&model.CCManageLog{
		LogType:      model.LogTypeAddLegion,
		TargetType:   "legion",
		TargetID:     legion.ID,
		Content:      "初始化创建军团",
		OperatorName: "系统初始化",
		CreateTime:   time.Now(),
	})

	db.Create(&model.CCFundLog{
		LogType:      model.FundLogTypeLegionRecharge,
		TargetType:   "legion",
		TargetID:     legion.ID,
		TargetName:   legion.LegionName,
		Amount:       10000000,
		BalanceAfter: 10000000,
		Content:      "初始化资金注入",
		OperatorName: "系统初始化",
		CreateTime:   time.Now(),
	})

	log.Println("CC data initialized successfully.")
}
