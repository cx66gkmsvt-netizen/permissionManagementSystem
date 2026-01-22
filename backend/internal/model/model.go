package model

import (
	"time"
)

// BaseModel 基础模型
type BaseModel struct {
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
}

// SysUser 用户表
type SysUser struct {
	UserID    int64      `json:"userId" gorm:"primaryKey;autoIncrement;column:user_id"`
	DeptID    *int64     `json:"deptId" gorm:"column:dept_id"`
	UserName  string     `json:"userName" gorm:"column:user_name;size:30;not null;uniqueIndex"`
	NickName  string     `json:"nickName" gorm:"column:nick_name;size:30"`
	Password  string     `json:"-" gorm:"column:password;size:100;not null"`
	Email     string     `json:"email" gorm:"column:email;size:50"`
	Phone     string     `json:"phone" gorm:"column:phone;size:11"`
	Avatar    string     `json:"avatar" gorm:"column:avatar;size:255"`
	Status    string     `json:"status" gorm:"column:status;size:1;default:0"` // 0正常 1停用
	DelFlag   string     `json:"-" gorm:"column:del_flag;size:1;default:0"`    // 0存在 2删除
	LoginIP   string     `json:"loginIp" gorm:"column:login_ip;size:128"`
	LoginDate *time.Time `json:"loginDate" gorm:"column:login_date"`
	CreateBy  *int64     `json:"createBy" gorm:"column:create_by"`
	BaseModel

	// 关联
	Dept  *SysDept  `json:"dept" gorm:"foreignKey:DeptID;references:DeptID"`
	Roles []SysRole `json:"roles" gorm:"many2many:sys_user_role;foreignKey:UserID;joinForeignKey:user_id;references:RoleID;joinReferences:role_id"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

// SysRole 角色表
type SysRole struct {
	RoleID    int64  `json:"roleId" gorm:"primaryKey;autoIncrement;column:role_id"`
	RoleName  string `json:"roleName" gorm:"column:role_name;size:30;not null"`
	RoleKey   string `json:"roleKey" gorm:"column:role_key;size:100;not null;uniqueIndex"`
	DataScope string `json:"dataScope" gorm:"column:data_scope;size:1;default:1"` // 1全部 2自定义 3本部门及以下 4本部门 5仅本人
	Sort      int    `json:"sort" gorm:"column:sort;default:0"`
	Status    string `json:"status" gorm:"column:status;size:1;default:0"`
	DelFlag   string `json:"-" gorm:"column:del_flag;size:1;default:0"`
	Remark    string `json:"remark" gorm:"column:remark;size:500"`
	BaseModel

	// 关联
	Menus []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;foreignKey:RoleID;joinForeignKey:role_id;references:MenuID;joinReferences:menu_id"`
	Depts []SysDept `json:"depts" gorm:"many2many:sys_role_dept;foreignKey:RoleID;joinForeignKey:role_id;references:DeptID;joinReferences:dept_id"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

// SysDept 部门表
type SysDept struct {
	DeptID    int64  `json:"deptId" gorm:"primaryKey;autoIncrement;column:dept_id"`
	ParentID  int64  `json:"parentId" gorm:"column:parent_id;default:0"`
	Ancestors string `json:"ancestors" gorm:"column:ancestors;size:500"` // 祖级列表
	DeptName  string `json:"deptName" gorm:"column:dept_name;size:30;not null"`
	Sort      int    `json:"sort" gorm:"column:sort;default:0"`
	Leader    string `json:"leader" gorm:"column:leader;size:20"`
	Phone     string `json:"phone" gorm:"column:phone;size:11"`
	Email     string `json:"email" gorm:"column:email;size:50"`
	Status    string `json:"status" gorm:"column:status;size:1;default:0"`
	DelFlag   string `json:"-" gorm:"column:del_flag;size:1;default:0"`
	BaseModel

	// 非数据库字段
	Children []*SysDept `json:"children" gorm:"-"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

// SysMenu 菜单表
type SysMenu struct {
	MenuID    int64  `json:"menuId" gorm:"primaryKey;autoIncrement;column:menu_id"`
	ParentID  int64  `json:"parentId" gorm:"column:parent_id;default:0"`
	MenuName  string `json:"menuName" gorm:"column:menu_name;size:50;not null"`
	MenuType  string `json:"menuType" gorm:"column:menu_type;size:1"` // M目录 C菜单 F按钮
	Path      string `json:"path" gorm:"column:path;size:200"`
	Component string `json:"component" gorm:"column:component;size:255"`
	Perms     string `json:"perms" gorm:"column:perms;size:100"` // 权限标识
	Icon      string `json:"icon" gorm:"column:icon;size:100"`
	Sort      int    `json:"sort" gorm:"column:sort;default:0"`
	Visible   string `json:"visible" gorm:"column:visible;size:1;default:0"` // 0显示 1隐藏
	Status    string `json:"status" gorm:"column:status;size:1;default:0"`
	BaseModel

	// 非数据库字段
	Children []*SysMenu `json:"children" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

// SysOperLog 操作日志表
type SysOperLog struct {
	OperID        int64     `json:"operId" gorm:"primaryKey;autoIncrement;column:oper_id"`
	Title         string    `json:"title" gorm:"column:title;size:50"`
	BusinessType  int       `json:"businessType" gorm:"column:business_type;default:0"` // 0其他 1新增 2修改 3删除
	Method        string    `json:"method" gorm:"column:method;size:100"`
	RequestMethod string    `json:"requestMethod" gorm:"column:request_method;size:10"`
	OperUserID    *int64    `json:"operUserId" gorm:"column:oper_user_id"`
	OperUserName  string    `json:"operUserName" gorm:"column:oper_user_name;size:50"`
	OperURL       string    `json:"operUrl" gorm:"column:oper_url;size:255"`
	OperIP        string    `json:"operIp" gorm:"column:oper_ip;size:128"`
	OperParam     string    `json:"operParam" gorm:"column:oper_param;type:text"`
	JSONResult    string    `json:"jsonResult" gorm:"column:json_result;type:text"`
	Status        string    `json:"status" gorm:"column:status;size:1;default:0"`
	ErrorMsg      string    `json:"errorMsg" gorm:"column:error_msg;type:text"`
	OperTime      time.Time `json:"operTime" gorm:"column:oper_time;autoCreateTime"`
}

func (SysOperLog) TableName() string {
	return "sys_oper_log"
}

// SysFollowUp 跟进记录表
type SysFollowUp struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id;comment:主键ID"`
	TargetType   string    `json:"targetType" gorm:"column:target_type;size:50;index;comment:目标类型"` // 目标类型，如 sys_user
	TargetID     int64     `json:"targetId" gorm:"column:target_id;index;comment:目标ID"`             // 目标ID
	Content      string    `json:"content" gorm:"column:content;type:text;comment:跟进内容"`            // 跟进内容
	OperUserID   *int64    `json:"operUserId" gorm:"column:oper_user_id;comment:操作人ID"`             // 操作人ID
	OperUserName string    `json:"operUserName" gorm:"column:oper_user_name;size:50;comment:操作人姓名"` // 操作人姓名
	OperTime     time.Time `json:"operTime" gorm:"column:oper_time;autoCreateTime;comment:操作时间"`    // 操作时间
	Remark       string    `json:"remark" gorm:"column:remark;size:500;comment:备注"`                 // 备注
}

func (SysFollowUp) TableName() string {
	return "sys_follow_up"
}
