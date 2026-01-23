package model

// CCMember CC成员表
type CCMember struct {
	ID           int64   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name         string  `json:"name" gorm:"column:name;size:50;not null"`                   // 老师名/姓名
	NickName     string  `json:"nickName" gorm:"column:nick_name;size:50"`                   // 昵称
	Mobile       string  `json:"mobile" gorm:"column:mobile;size:11;not null"`               // 手机号/账号
	WeChat       string  `json:"wechat" gorm:"column:wechat;size:50"`                        // 微信号
	Cno          string  `json:"cno" gorm:"column:cno;size:20"`                              // 座席号
	CloudAccount string  `json:"cloudAccount" gorm:"column:cloud_account;size:50"`           // 云客账号
	TeamID       int64   `json:"teamId" gorm:"column:team_id"`                               // 所属军团ID
	SquadID      int64   `json:"squadId" gorm:"column:squad_id"`                             // 所属战队ID
	Balance      float64 `json:"balance" gorm:"column:balance;type:decimal(10,2);default:0"` // 个人资金余额
	Status       string  `json:"status" gorm:"column:status;size:1;default:0"`               // 0正常 1停用
	DelFlag      string  `json:"-" gorm:"column:del_flag;size:1;default:0"`                  // 0存在 2删除
	CreateBy     string  `json:"createBy" gorm:"column:create_by;size:50"`
	UpdateBy     string  `json:"updateBy" gorm:"column:update_by;size:50"`
	BaseModel

	// 关联字段 (View Only)
	TeamName  string `json:"teamName" gorm:"-"`
	SquadName string `json:"squadName" gorm:"-"`
}

func (CCMember) TableName() string {
	return "cc_member"
}

// CCQuery 查询参数
type CCQuery struct {
	PageQuery
	Name    string `form:"name"`
	Mobile  string `form:"mobile"`
	TeamID  *int64 `form:"teamId"`
	SquadID *int64 `form:"squadId"`
	Status  string `form:"status"`
}
