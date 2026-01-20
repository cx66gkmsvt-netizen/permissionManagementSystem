package model

// LoginRequest 登录请求
type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
}

// UserInfo 用户信息响应
type UserInfo struct {
	User        *SysUser `json:"user"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// PageQuery 分页查询
type PageQuery struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func (p *PageQuery) GetOffset() int {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.PageNum - 1) * p.PageSize
}

// PageResult 分页结果
type PageResult struct {
	Total int64       `json:"total"`
	Rows  interface{} `json:"rows"`
}

// UserQuery 用户查询
type UserQuery struct {
	PageQuery
	UserName string `json:"userName" form:"userName"`
	Phone    string `json:"phone" form:"phone"`
	Status   string `json:"status" form:"status"`
	DeptID   *int64 `json:"deptId" form:"deptId"`
}

// RoleQuery 角色查询
type RoleQuery struct {
	PageQuery
	RoleName string `json:"roleName" form:"roleName"`
	RoleKey  string `json:"roleKey" form:"roleKey"`
	Status   string `json:"status" form:"status"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	DeptID   *int64  `json:"deptId"`
	UserName string  `json:"userName" binding:"required"`
	NickName string  `json:"nickName"`
	Password string  `json:"password" binding:"required,min=6,max=20"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Status   string  `json:"status"`
	RoleIDs  []int64 `json:"roleIds"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	DeptID   *int64  `json:"deptId"`
	NickName string  `json:"nickName"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Status   string  `json:"status"`
	RoleIDs  []int64 `json:"roleIds"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName  string  `json:"roleName" binding:"required"`
	RoleKey   string  `json:"roleKey" binding:"required"`
	DataScope string  `json:"dataScope"`
	Sort      int     `json:"sort"`
	Status    string  `json:"status"`
	Remark    string  `json:"remark"`
	MenuIDs   []int64 `json:"menuIds"`
	DeptIDs   []int64 `json:"deptIds"` // 自定义数据权限
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName  string  `json:"roleName"`
	RoleKey   string  `json:"roleKey"`
	DataScope string  `json:"dataScope"`
	Sort      int     `json:"sort"`
	Status    string  `json:"status"`
	Remark    string  `json:"remark"`
	MenuIDs   []int64 `json:"menuIds"`
	DeptIDs   []int64 `json:"deptIds"`
}

// CreateDeptRequest 创建部门请求
type CreateDeptRequest struct {
	ParentID int64  `json:"parentId"`
	DeptName string `json:"deptName" binding:"required"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID  int64  `json:"parentId"`
	MenuName  string `json:"menuName" binding:"required"`
	MenuType  string `json:"menuType" binding:"required"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Perms     string `json:"perms"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	Visible   string `json:"visible"`
	Status    string `json:"status"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
}
