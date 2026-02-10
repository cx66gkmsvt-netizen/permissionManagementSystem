package handler

import (
	"user-center/internal/model"
	"user-center/internal/repository"

	"github.com/gin-gonic/gin"
)

// DataScope 数据范围信息
type DataScope struct {
	IsAdmin  bool   // 是否管理员（无限制）
	IsCCUser bool   // 是否CC用户
	RoleType string // CC角色类型
	CCID     int64  // 当前CC的ID
	SquadID  *int64 // 所属战队ID
	TeamID   *int64 // 所属团队ID
	LegionID *int64 // 所属军团ID
}

// getDataScope 获取当前用户的数据范围
// 根据登录用户的角色判断其数据可见范围
func getDataScope(c *gin.Context) *DataScope {
	userID := c.GetInt64("userID")

	// 超管(userID==1)不做限制
	if userID == 1 {
		return &DataScope{IsAdmin: true}
	}

	// 检查用户是否有CC角色
	userRepo := repository.NewUserRepository()
	roles, _ := userRepo.GetUserRoles(userID)

	var ccRoleKey string
	for _, role := range roles {
		switch role.RoleKey {
		case "cc":
			ccRoleKey = model.RoleTypeCC
		case "cc_squad_leader":
			ccRoleKey = model.RoleTypeSquadLeader
		case "cc_team_leader":
			ccRoleKey = model.RoleTypeTeamLeader
		case "cc_legion_leader":
			ccRoleKey = model.RoleTypeLegionLeader
		case "admin":
			return &DataScope{IsAdmin: true}
		}
	}

	// 非CC角色用户，视为管理员权限
	if ccRoleKey == "" {
		return &DataScope{IsAdmin: true}
	}

	// 获取CC成员信息
	ccRepo := repository.NewCCRepository()
	cc, err := ccRepo.Get(userID)
	if err != nil {
		// CC记录不存在，按管理员处理
		return &DataScope{IsAdmin: true}
	}

	return &DataScope{
		IsAdmin:  false,
		IsCCUser: true,
		RoleType: ccRoleKey,
		CCID:     cc.ID,
		SquadID:  cc.SquadID,
		TeamID:   cc.TeamID,
		LegionID: cc.LegionID,
	}
}

// applyDataScopeToCCQuery 将数据范围应用到CC查询条件
func applyDataScopeToCCQuery(query *model.CCQuery, scope *DataScope) {
	if scope.IsAdmin {
		return // 管理员不做限制
	}

	switch scope.RoleType {
	case model.RoleTypeCC:
		// 普通CC只能看自己
		query.CCID = &scope.CCID
	case model.RoleTypeSquadLeader:
		// 战队长看自己战队下的人
		if scope.SquadID != nil {
			query.SquadID = scope.SquadID
		} else {
			query.CCID = &scope.CCID // fallback: 只看自己
		}
	case model.RoleTypeTeamLeader:
		// 团长看自己团队下的人
		if scope.TeamID != nil {
			query.TeamID = scope.TeamID
		} else {
			query.CCID = &scope.CCID
		}
	case model.RoleTypeLegionLeader:
		// 军团长看自己军团下的人
		if scope.LegionID != nil {
			query.LegionID = scope.LegionID
		} else {
			query.CCID = &scope.CCID
		}
	}
}

// applyDataScopeToLeadAllocationQuery 将数据范围应用到例子分配查询条件
func applyDataScopeToLeadAllocationQuery(query *model.LeadAllocationQuery, scope *DataScope) {
	if scope.IsAdmin {
		return
	}

	switch scope.RoleType {
	case model.RoleTypeCC:
		query.CCID = &scope.CCID
	case model.RoleTypeSquadLeader:
		if scope.SquadID != nil {
			query.SquadID = scope.SquadID
		} else {
			query.CCID = &scope.CCID
		}
	case model.RoleTypeTeamLeader:
		if scope.TeamID != nil {
			query.TeamID = scope.TeamID
		} else {
			query.CCID = &scope.CCID
		}
	case model.RoleTypeLegionLeader:
		if scope.LegionID != nil {
			query.LegionID = scope.LegionID
		} else {
			query.CCID = &scope.CCID
		}
	}
}
