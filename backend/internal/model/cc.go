package model

// CCMember CC成员表
type CCMember struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:CC成员ID"`
	Name     string `json:"name" gorm:"column:name;size:50;not null;comment:老师名/姓名"`
	NickName string `json:"nickName" gorm:"column:nick_name;size:50;comment:昵称"`
	Mobile   string `json:"mobile" gorm:"column:mobile;size:11;not null;uniqueIndex;comment:手机号/账号"`
	WeChat   string `json:"wechat" gorm:"column:wechat;size:50;comment:微信号"`

	// 角色类型
	RoleType string `json:"roleType" gorm:"column:role_type;size:20;default:cc;comment:角色类型(cc/squad_leader/team_leader/legion_leader)"`

	// 座席号配置
	Cno1          string `json:"cno1" gorm:"column:cno1;size:20;comment:天润CNO1"`
	Cno2          string `json:"cno2" gorm:"column:cno2;size:20;comment:天润CNO2"`
	CloudAccount1 string `json:"cloudAccount1" gorm:"column:cloud_account1;size:50;comment:云客账号1"`
	CloudAccount2 string `json:"cloudAccount2" gorm:"column:cloud_account2;size:50;comment:云客账号2"`
	RonglianSeat  string `json:"ronglianSeat" gorm:"column:ronglian_seat;size:50;comment:容联座席号"`
	DiankongSeat  string `json:"diankongSeat" gorm:"column:diankong_seat;size:50;comment:点控云座席号"`
	HeliAccount   string `json:"heliAccount" gorm:"column:heli_account;size:50;comment:合力亿捷账号"`
	BaichuanSeat  string `json:"baichuanSeat" gorm:"column:baichuan_seat;size:50;comment:百川智通座席号"`

	// 智能外呼配置
	DiankongOutboundPool int `json:"diankongOutboundPool" gorm:"column:diankong_outbound_pool;default:1;comment:点控云公海智能外呼(0关闭 1开启)"`
	DiankongOutboundList int `json:"diankongOutboundList" gorm:"column:diankong_outbound_list;default:0;comment:点控云客户列表智能外呼(0关闭 1开启)"`
	BaichuanOutboundPool int `json:"baichuanOutboundPool" gorm:"column:baichuan_outbound_pool;default:0;comment:百川公海智能外呼(0关闭 1开启)"`
	BaichuanOutboundList int `json:"baichuanOutboundList" gorm:"column:baichuan_outbound_list;default:0;comment:百川客户列表智能外呼(0关闭 1开启)"`

	// 组织架构
	SquadID  *int64 `json:"squadId" gorm:"column:squad_id;comment:所属战队ID"`
	TeamID   *int64 `json:"teamId" gorm:"column:team_id;comment:所属团队ID"`
	LegionID *int64 `json:"legionId" gorm:"column:legion_id;comment:所属军团ID"`

	// 资金和业绩
	Balance            int64 `json:"balance" gorm:"column:balance;default:0;comment:个人余额(分)"`
	MonthlyPerformance int64 `json:"monthlyPerformance" gorm:"column:monthly_performance;default:0;comment:当月业绩(分)"`
	PerformanceRank    int   `json:"performanceRank" gorm:"column:performance_rank;default:0;comment:业绩排名"`

	// 在班状态
	AttendanceStatus string `json:"attendanceStatus" gorm:"column:attendance_status;size:1;default:2;comment:在班状态(1在班 2休班 3请假)"`

	// 通时通次目标
	CallDurationTarget int `json:"callDurationTarget" gorm:"column:call_duration_target;default:0;comment:通时目标(秒)"`
	CallCountTarget    int `json:"callCountTarget" gorm:"column:call_count_target;default:0;comment:通次目标"`

	// 状态标志
	Status    string `json:"status" gorm:"column:status;size:1;default:0;comment:状态(0正常 1停用)"`
	IsBlocked string `json:"isBlocked" gorm:"column:is_blocked;size:1;default:0;comment:是否屏蔽(0否 1是)"`
	DelFlag   string `json:"-" gorm:"column:del_flag;size:1;default:0;comment:删除标志(0存在 2删除)"`

	CreateBy string `json:"createBy" gorm:"column:create_by;size:50;comment:创建者"`
	UpdateBy string `json:"updateBy" gorm:"column:update_by;size:50;comment:更新者"`
	BaseModel

	// 关联字段 (View Only) - 前端显示用
	SquadName        string  `json:"squadName" gorm:"-"`
	SquadLeaderName  string  `json:"squadLeaderName" gorm:"-"`
	TeamName         string  `json:"teamName" gorm:"-"`
	TeamLeaderName   string  `json:"teamLeaderName" gorm:"-"`
	LegionName       string  `json:"legionName" gorm:"-"`
	LegionLeaderName string  `json:"legionLeaderName" gorm:"-"`
	BalanceYuan      float64 `json:"balanceYuan" gorm:"-"`     // 余额(元)，前端显示用
	PerformanceYuan  float64 `json:"performanceYuan" gorm:"-"` // 业绩(元)，前端显示用
}

func (CCMember) TableName() string {
	return "cc_member"
}

// 角色类型常量
const (
	RoleTypeCC           = "cc"            // 普通CC
	RoleTypeSquadLeader  = "squad_leader"  // 战队长
	RoleTypeTeamLeader   = "team_leader"   // 团长
	RoleTypeLegionLeader = "legion_leader" // 军团长
)

// CCQuery 查询参数
type CCQuery struct {
	PageQuery
	CCID             *int64 `form:"ccId"`
	Name             string `form:"name"`
	NickName         string `form:"nickName"`
	Mobile           string `form:"mobile"`
	RoleType         string `form:"roleType"`
	SquadID          *int64 `form:"squadId"`
	TeamID           *int64 `form:"teamId"`
	LegionID         *int64 `form:"legionId"`
	AttendanceStatus string `form:"attendanceStatus"`
	Status           string `form:"status"`
	IsBlocked        string `form:"isBlocked"`
}

// CCUpdateDTO CC更新DTO
type CCUpdateDTO struct {
	Name                 string `json:"name"`
	NickName             string `json:"nickName"`
	WeChat               string `json:"wechat"`
	Cno1                 string `json:"cno1"`
	Cno2                 string `json:"cno2"`
	CloudAccount1        string `json:"cloudAccount1"`
	CloudAccount2        string `json:"cloudAccount2"`
	RonglianSeat         string `json:"ronglianSeat"`
	DiankongSeat         string `json:"diankongSeat"`
	HeliAccount          string `json:"heliAccount"`
	BaichuanSeat         string `json:"baichuanSeat"`
	DiankongOutboundPool *int   `json:"diankongOutboundPool"`
	DiankongOutboundList *int   `json:"diankongOutboundList"`
	BaichuanOutboundPool *int   `json:"baichuanOutboundPool"`
	BaichuanOutboundList *int   `json:"baichuanOutboundList"`
	SquadID              *int64 `json:"squadId"`
	StructAdjustDate     string `json:"structAdjustDate"`
	TransactionAmount    *int64 `json:"transactionAmount"`
}
