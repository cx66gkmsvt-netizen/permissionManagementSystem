package model

import (
	"time"
)

// CCLegion 军团表
type CCLegion struct {
	ID         int64  `json:"id,string" gorm:"primaryKey;autoIncrement;column:id;comment:军团ID"`
	LegionName string `json:"legionName" gorm:"column:legion_name;size:50;not null;uniqueIndex;comment:军团名称"`
	LeaderID   *int64 `json:"leaderId,string" gorm:"column:leader_id;comment:军团长ID"`
	Balance    int64  `json:"balance" gorm:"column:balance;default:0;comment:军团余额(分)"`
	Status     string `json:"status" gorm:"column:status;size:1;default:0;comment:状态(0正常 1停用)"`
	DelFlag    string `json:"-" gorm:"column:del_flag;size:1;default:0;comment:删除标志(0存在 2删除)"`
	CreateBy   string `json:"createBy" gorm:"column:create_by;size:50;comment:创建者"`
	UpdateBy   string `json:"updateBy" gorm:"column:update_by;size:50;comment:更新者"`
	BaseModel

	// 关联字段 (View Only)
	LeaderName         string  `json:"leaderName" gorm:"->"`
	MonthlyPerformance float64 `json:"monthlyPerformance" gorm:"->"` // 当月业绩(元)
	PerformanceRank    int     `json:"performanceRank" gorm:"->"`    // 业绩排名
}

func (CCLegion) TableName() string {
	return "cc_legion"
}

// CCTeam 团队表
type CCTeam struct {
	ID           int64  `json:"id,string" gorm:"primaryKey;autoIncrement;column:id;comment:团队ID"`
	TeamName     string `json:"teamName" gorm:"column:team_name;size:50;not null;comment:团队名称"`
	BusinessType string `json:"businessType" gorm:"column:business_type;size:20;comment:业务类型"`
	LeaderID     *int64 `json:"leaderId,string" gorm:"column:leader_id;comment:团长ID"`
	LegionID     *int64 `json:"legionId,string" gorm:"column:legion_id;comment:所属军团ID"`
	Balance      int64  `json:"balance" gorm:"column:balance;default:0;comment:团队余额(分)"`
	Status       string `json:"status" gorm:"column:status;size:1;default:0;comment:状态(0正常 1停用)"`
	DelFlag      string `json:"-" gorm:"column:del_flag;size:1;default:0;comment:删除标志(0存在 2删除)"`
	CreateBy     string `json:"createBy" gorm:"column:create_by;size:50;comment:创建者"`
	UpdateBy     string `json:"updateBy" gorm:"column:update_by;size:50;comment:更新者"`
	BaseModel

	// 关联字段 (View Only)
	LeaderName         string  `json:"leaderName" gorm:"->"`
	LegionName         string  `json:"legionName" gorm:"->"`
	LegionLeaderName   string  `json:"legionLeaderName" gorm:"->"`
	MonthlyPerformance float64 `json:"monthlyPerformance" gorm:"->"` // 当月业绩(元)
	PerformanceRank    int     `json:"performanceRank" gorm:"->"`    // 业绩排名
}

func (CCTeam) TableName() string {
	return "cc_team"
}

// CCSquad 战队表
type CCSquad struct {
	ID        int64  `json:"id,string" gorm:"primaryKey;autoIncrement;column:id;comment:战队ID"`
	SquadName string `json:"squadName" gorm:"column:squad_name;size:50;not null;uniqueIndex;comment:战队名称"`
	LeaderID  *int64 `json:"leaderId,string" gorm:"column:leader_id;comment:战队长ID"`
	TeamID    int64  `json:"teamId,string" gorm:"column:team_id;not null;comment:所属团队ID"`
	Balance   int64  `json:"balance" gorm:"column:balance;default:0;comment:战队余额(分)"`
	Status    string `json:"status" gorm:"column:status;size:1;default:0;comment:状态(0正常 1停用)"`
	DelFlag   string `json:"-" gorm:"column:del_flag;size:1;default:0;comment:删除标志(0存在 2删除)"`
	CreateBy  string `json:"createBy" gorm:"column:create_by;size:50;comment:创建者"`
	UpdateBy  string `json:"updateBy" gorm:"column:update_by;size:50;comment:更新者"`
	BaseModel

	// 关联字段 (View Only)
	LeaderName         string  `json:"leaderName" gorm:"->"`
	TeamName           string  `json:"teamName" gorm:"->"`
	TeamLeaderName     string  `json:"teamLeaderName" gorm:"->"`
	LegionID           *int64  `json:"legionId,string" gorm:"->"`
	LegionName         string  `json:"legionName" gorm:"->"`
	LegionLeaderName   string  `json:"legionLeaderName" gorm:"->"`
	MonthlyPerformance float64 `json:"monthlyPerformance" gorm:"->"` // 当月业绩(元)
	PerformanceRank    int     `json:"performanceRank" gorm:"->"`    // 业绩排名
}

func (CCSquad) TableName() string {
	return "cc_squad"
}

// CCManageLog CC管理日志表
type CCManageLog struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:日志ID"`
	LogType      string    `json:"logType" gorm:"column:log_type;size:50;not null;comment:日志类型"`
	TargetType   string    `json:"targetType" gorm:"column:target_type;size:20;comment:目标类型(legion/team/squad/cc)"`
	TargetID     int64     `json:"targetId" gorm:"column:target_id;comment:目标ID"`
	Content      string    `json:"content" gorm:"column:content;type:text;comment:日志内容"`
	OperatorID   int64     `json:"operatorId" gorm:"column:operator_id;comment:操作人ID"`
	OperatorName string    `json:"operatorName" gorm:"column:operator_name;size:50;comment:操作人名称"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
}

func (CCManageLog) TableName() string {
	return "cc_manage_log"
}

// 日志类型常量
const (
	LogTypeAddLegion         = "add_legion"          // 添加军团
	LogTypeModifyLegionInfo  = "modify_legion_info"  // 修改军团信息
	LogTypeAddTeam           = "add_team"            // 添加团队
	LogTypeModifyTeamInfo    = "modify_team_info"    // 修改团队信息
	LogTypeModifyTeamStruct  = "modify_team_struct"  // 修改团队结构
	LogTypeAddSquad          = "add_squad"           // 添加战队
	LogTypeModifySquadInfo   = "modify_squad_info"   // 修改战队信息
	LogTypeModifySquadStruct = "modify_squad_struct" // 修改战队结构
	LogTypeModifyCCInfo      = "modify_cc_info"      // 修改CC信息
	LogTypeModifyCCStruct    = "modify_cc_struct"    // 修改CC结构
)

// CCFundLog CC资金日志表
type CCFundLog struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:日志ID"`
	LogType         string    `json:"logType" gorm:"column:log_type;size:50;not null;comment:日志类型"`
	TargetType      string    `json:"targetType" gorm:"column:target_type;size:20;not null;comment:目标类型(legion/team/squad/cc)"`
	TargetID        int64     `json:"targetId" gorm:"column:target_id;not null;comment:目标ID"`
	TargetName      string    `json:"targetName" gorm:"column:target_name;size:50;comment:目标名称"`
	Amount          int64     `json:"amount" gorm:"column:amount;default:0;comment:变动金额(分)"`
	BalanceBefore   int64     `json:"balanceBefore" gorm:"column:balance_before;default:0;comment:变动前余额(分)"`
	BalanceAfter    int64     `json:"balanceAfter" gorm:"column:balance_after;default:0;comment:变动后余额(分)"`
	RelatedCCID     *int64    `json:"relatedCcId" gorm:"column:related_cc_id;comment:关联CC ID"`
	RelatedSquadID  *int64    `json:"relatedSquadId" gorm:"column:related_squad_id;comment:关联战队ID"`
	RelatedTeamID   *int64    `json:"relatedTeamId" gorm:"column:related_team_id;comment:关联团队ID"`
	RelatedLegionID *int64    `json:"relatedLegionId" gorm:"column:related_legion_id;comment:关联军团ID"`
	Content         string    `json:"content" gorm:"column:content;type:text;comment:日志内容详情"`
	Reason          string    `json:"reason" gorm:"column:reason;size:500;comment:变动原因"`
	OperatorID      int64     `json:"operatorId" gorm:"column:operator_id;comment:操作人ID"`
	OperatorName    string    `json:"operatorName" gorm:"column:operator_name;size:50;comment:操作人名称"`
	CreateTime      time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
}

func (CCFundLog) TableName() string {
	return "cc_fund_log"
}

// 资金日志类型常量
const (
	FundLogTypePerformanceIncrease = "performance_increase"  // 业绩流水增加
	FundLogTypePerformanceDecrease = "performance_decrease"  // 业绩流水减少
	FundLogTypePromoteSquadLeader  = "promote_squad_leader"  // 晋升战队长
	FundLogTypePromoteTeamLeader   = "promote_team_leader"   // 晋升团长
	FundLogTypePromoteLegionLeader = "promote_legion_leader" // 晋升军团长
	FundLogTypeModifyCCStruct      = "modify_cc_struct"      // 修改CC结构
	FundLogTypeAddTeam             = "add_team"              // 添加团队
	FundLogTypeAddSquad            = "add_squad"             // 添加战队
	FundLogTypeCCBalanceEdit       = "cc_balance_edit"       // CC余额修改
	FundLogTypeSquadBalanceEdit    = "squad_balance_edit"    // 战队余额修改
	FundLogTypeTeamBalanceEdit     = "team_balance_edit"     // 团队余额修改
	FundLogTypeLegionBalanceEdit   = "legion_balance_edit"   // 军团余额修改
	FundLogTypeCCReward            = "cc_reward"             // CC奖励&扣除
	FundLogTypeSquadReward         = "squad_reward"          // 战队奖励&扣除
	FundLogTypeTeamReward          = "team_reward"           // 团队奖励&扣除
	FundLogTypeLegionReward        = "legion_reward"         // 军团奖励&扣除
	FundLogTypeLegionRecharge      = "legion_recharge"       // 军团充值
	FundLogTypeLegionTransfer      = "legion_transfer"       // 军团转账
	FundLogTypeTeamRecharge        = "team_recharge"         // 团队充值
	FundLogTypeTeamTransfer        = "team_transfer"         // 团队转账
	FundLogTypeSquadRecharge       = "squad_recharge"        // 战队充值
	FundLogTypeSquadTransfer       = "squad_transfer"        // 战队转账
	FundLogTypeAdminTransfer       = "admin_transfer"        // 管理员转账
)

// CCAttendance 在班记录表
type CCAttendance struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:记录ID"`
	CCID           int64     `json:"ccId" gorm:"column:cc_id;not null;comment:CC成员ID"`
	AttendanceDate time.Time `json:"attendanceDate" gorm:"column:attendance_date;type:date;not null;comment:在班日期"`
	Status         string    `json:"status" gorm:"column:status;size:1;not null;comment:在班状态(1在班 2休班 3请假)"`
	OperatorID     *int64    `json:"operatorId" gorm:"column:operator_id;comment:操作人ID"`
	CreateTime     time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`

	// 关联字段 (View Only)
	CCName     string `json:"ccName" gorm:"->"`
	LegionName string `json:"legionName" gorm:"->"`
	TeamName   string `json:"teamName" gorm:"->"`
	SquadName  string `json:"squadName" gorm:"->"`
}

func (CCAttendance) TableName() string {
	return "cc_attendance"
}

// 在班状态常量
const (
	AttendanceStatusOnDuty  = "1" // 在班
	AttendanceStatusOffDuty = "2" // 休班
	AttendanceStatusLeave   = "3" // 请假
)

// CCLeadAllocation 例子分配表
type CCLeadAllocation struct {
	ID                 int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:记录ID"`
	CCID               int64     `json:"ccId" gorm:"column:cc_id;not null;comment:CC成员ID"`
	AllocationDate     time.Time `json:"allocationDate" gorm:"column:allocation_date;type:date;not null;comment:分配日期"`
	ExpectedAllocation int       `json:"expectedAllocation" gorm:"column:expected_allocation;default:0;comment:预计分配"`
	ActualAllocation   int       `json:"actualAllocation" gorm:"column:actual_allocation;default:0;comment:实际分配"`
	ExpectedSupplement int       `json:"expectedSupplement" gorm:"column:expected_supplement;default:0;comment:预计补发"`
	ActualSupplement   int       `json:"actualSupplement" gorm:"column:actual_supplement;default:0;comment:实际补发"`
	Overdraft          int       `json:"overdraft" gorm:"column:overdraft;default:0;comment:透支数量"`
	ProcessedOverdraft int       `json:"processedOverdraft" gorm:"column:processed_overdraft;default:0;comment:已处理透支"`
	PendingOverdraft   int       `json:"pendingOverdraft" gorm:"column:pending_overdraft;default:0;comment:待处理透支"`
	IsAllocated        string    `json:"isAllocated" gorm:"column:is_allocated;size:1;default:0;comment:是否分配例子(0否 1是)"`
	AllocationRule     string    `json:"allocationRule" gorm:"column:allocation_rule;size:10;comment:分配规则(A节假日 B工作日无补偿 C工作日)"`
	AllocationReason   string    `json:"allocationReason" gorm:"column:allocation_reason;type:text;comment:分配/未分配原因"`
	CreateTime         time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime         time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`

	// 关联字段 (View Only)
	CCName     string `json:"ccName" gorm:"->"`
	LegionName string `json:"legionName" gorm:"->"`
	TeamName   string `json:"teamName" gorm:"->"`
	SquadName  string `json:"squadName" gorm:"->"`
}

func (CCLeadAllocation) TableName() string {
	return "cc_lead_allocation"
}

// 分配规则常量
const (
	AllocationRuleHoliday       = "A" // 节假日
	AllocationRuleWorkdayNoComp = "B" // 工作日（无例子补偿）
	AllocationRuleWorkday       = "C" // 工作日
)

// CCFundConfig 资金配置表
type CCFundConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:配置ID"`
	ConfigType  string    `json:"configType" gorm:"column:config_type;size:50;not null;comment:配置类型"`
	RankOrID    string    `json:"rankOrId" gorm:"column:rank_or_id;size:50;comment:ID/排名序号"`
	Amount      int64     `json:"amount" gorm:"column:amount;default:0;comment:金额(分)"`
	ConfigMonth string    `json:"configMonth" gorm:"column:config_month;size:7;comment:配置月份(YYYY-MM)"`
	CreateBy    string    `json:"createBy" gorm:"column:create_by;size:50;comment:创建者"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
}

func (CCFundConfig) TableName() string {
	return "cc_fund_config"
}

// 资金配置类型常量
const (
	FundConfigPersonalSubsidy = "personal_subsidy" // 个人补贴
	FundConfigLeadCost        = "lead_cost"        // 例子成本
	FundConfigSquadSubsidy    = "squad_subsidy"    // 战队补贴
	FundConfigTeamSubsidy     = "team_subsidy"     // 团队补贴
	FundConfigPersonalReward  = "personal_reward"  // 个人奖励
	FundConfigSquadReward     = "squad_reward"     // 战队奖励
	FundConfigTeamReward      = "team_reward"      // 团队奖励
	FundConfigSquadLicense    = "squad_license"    // 战队许可
	FundConfigTeamLicense     = "team_license"     // 团队许可
)

// Query structs

// LegionQuery 军团查询参数
type LegionQuery struct {
	PageQuery
	LegionName string `form:"legionName"`
}

// TeamQuery 团队查询参数
type TeamQuery struct {
	PageQuery
	TeamName     string `form:"teamName"`
	BusinessType string `form:"businessType"`
	LegionID     *int64 `form:"legionId"`
}

// SquadQuery 战队查询参数
type SquadQuery struct {
	PageQuery
	SquadName string `form:"squadName"`
	TeamID    *int64 `form:"teamId"`
	LegionID  *int64 `form:"legionId"`
}

// FundLogQuery 资金日志查询参数
type FundLogQuery struct {
	PageQuery
	TargetType string `form:"targetType"`
	TargetID   int64  `form:"targetId"`
	LogType    string `form:"logType"`
	BillType   string `form:"billType"` // flow=流水, non_flow=非流水, all=全部
}

// AttendanceQuery 在班记录查询参数
type AttendanceQuery struct {
	Dates []string `form:"dates"` // 日期列表
}

// LeadAllocationQuery 例子分配查询参数
type LeadAllocationQuery struct {
	PageQuery
	CCID             *int64 `form:"ccId"`
	CCName           string `form:"ccName"`
	LegionID         *int64 `form:"legionId"`
	TeamID           *int64 `form:"teamId"`
	SquadID          *int64 `form:"squadId"`
	AttendanceStatus string `form:"attendanceStatus"`
}

// DTO structures

// FundEditDTO 编辑余额DTO
type FundEditDTO struct {
	EditType string `json:"editType" binding:"required,oneof=set adjust"` // set=设置为, adjust=增减
	Amount   int64  `json:"amount" binding:"required"`                    // 金额(分)
	Reason   string `json:"reason" binding:"required"`                    // 修改原因
}

// FundRechargeDTO 充值DTO
type FundRechargeDTO struct {
	Amount int64 `json:"amount" binding:"required,gt=0"` // 充值金额(分)
}

// FundTransferDTO 转账DTO
type FundTransferDTO struct {
	Amount      int64 `json:"amount" binding:"required,gt=0"` // 转账金额(分)
	RecipientID int64 `json:"recipientId" binding:"required"` // 收款人CC ID
}

// AdminTransferDTO 管理员转账DTO
type AdminTransferDTO struct {
	FromType string `json:"fromType" binding:"required,oneof=legion team squad cc"` // 转出方类型
	FromID   int64  `json:"fromId" binding:"required"`                              // 转出方ID
	ToType   string `json:"toType" binding:"required,oneof=legion team squad cc"`   // 转入方类型
	ToID     int64  `json:"toId" binding:"required"`                                // 转入方ID
	Amount   int64  `json:"amount" binding:"required,gt=0"`                         // 转账金额(分)
}

// TransactionConfirmDTO 交易确认DTO
type TransactionConfirmDTO struct {
	Amount int64 `json:"amount" binding:"required,gt=0"` // 交易金额(分)
}

// LegionCreateDTO 创建军团DTO
type LegionCreateDTO struct {
	LegionName string `json:"legionName" binding:"required,max=10"` // 军团名称
}

// LegionUpdateDTO 更新军团DTO
type LegionUpdateDTO struct {
	LegionName        string `json:"legionName" binding:"required,max=10"` // 军团名称
	LeaderID          *int64 `json:"leaderId,string"`                      // 军团长ID
	TransactionAmount *int64 `json:"transactionAmount"`                    // 交易金额(分)，晋升军团长时需要
}

// TeamCreateDTO 创建团队DTO
type TeamCreateDTO struct {
	TeamName          string `json:"teamName" binding:"required,max=10"` // 团队名称
	BusinessType      string `json:"businessType" binding:"required"`    // 业务类型
	LegionID          *int64 `json:"legionId,string"`                    // 所属军团ID
	StructAdjustDate  string `json:"structAdjustDate"`                   // 架构调整时间
	TransactionAmount int64  `json:"transactionAmount"`                  // 交易金额(分)
}

// TeamUpdateDTO 更新团队DTO
type TeamUpdateDTO struct {
	TeamName          string `json:"teamName" binding:"required,max=10"` // 团队名称
	BusinessType      string `json:"businessType" binding:"required"`    // 业务类型
	LeaderID          *int64 `json:"leaderId,string"`                    // 团长ID
	LegionID          *int64 `json:"legionId,string"`                    // 所属军团ID
	LeaderAdjustDate  string `json:"leaderAdjustDate"`                   // 团长调整时间
	StructAdjustDate  string `json:"structAdjustDate"`                   // 架构调整时间
	TransactionAmount *int64 `json:"transactionAmount"`                  // 交易金额(分)
}

// SquadCreateDTO 创建战队DTO
type SquadCreateDTO struct {
	SquadName         string `json:"squadName" binding:"required,max=10"` // 战队名称
	TeamID            int64  `json:"teamId,string" binding:"required"`    // 所属团队ID
	StructAdjustDate  string `json:"structAdjustDate" binding:"required"` // 架构调整时间
	TransactionAmount int64  `json:"transactionAmount"`                   // 交易金额(分)
}

// SquadUpdateDTO 更新战队DTO
type SquadUpdateDTO struct {
	SquadName         string `json:"squadName" binding:"required,max=10"` // 战队名称
	TeamID            int64  `json:"teamId,string" binding:"required"`    // 所属团队ID
	LeaderID          *int64 `json:"leaderId,string"`                     // 战队长ID
	LeaderAdjustDate  string `json:"leaderAdjustDate"`                    // 战队长调整时间
	StructAdjustDate  string `json:"structAdjustDate"`                    // 架构调整时间
	TransactionAmount *int64 `json:"transactionAmount"`                   // 交易金额(分)
}

// AttendanceUpdateDTO 更新在班状态DTO
type AttendanceUpdateDTO struct {
	Status string `json:"status" binding:"required,oneof=1 2 3"` // 在班状态
}

// AttendanceBatchUpdateItem 批量更新在班状态项
type AttendanceBatchUpdateItem struct {
	CCID   int64  `json:"ccId" binding:"required"`
	Date   string `json:"date" binding:"required"`
	Status string `json:"status" binding:"required,oneof=1 2 3"`
}

// AttendanceBatchUpdateDTO 批量更新在班状态DTO
type AttendanceBatchUpdateDTO struct {
	Items []AttendanceBatchUpdateItem `json:"items" binding:"required"`
}

// LeadAllocationUpdateDTO 例子分配更新DTO
type LeadAllocationUpdateDTO struct {
	ExpectedAllocation int    `json:"expectedAllocation"`
	ActualAllocation   int    `json:"actualAllocation"`
	ExpectedSupplement int    `json:"expectedSupplement"`
	ActualSupplement   int    `json:"actualSupplement"`
	Overdraft          int    `json:"overdraft"`
	ProcessedOverdraft int    `json:"processedOverdraft"`
	PendingOverdraft   int    `json:"pendingOverdraft"`
	IsAllocated        string `json:"isAllocated"`
	AllocationRule     string `json:"allocationRule"`
	AllocationReason   string `json:"allocationReason"`
}

// LeadAllocationBatchUpdateDTO 批量更新是否分配DTO
type LeadAllocationBatchUpdateDTO struct {
	CCIDs       []int64 `json:"ccIds" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	IsAllocated string  `json:"isAllocated" binding:"required,oneof=0 1"`
	Reason      string  `json:"reason"`
}
