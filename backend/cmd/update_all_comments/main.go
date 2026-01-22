package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DSN = "root:Zhixin123..@tcp(43.138.1.141:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	fmt.Println("开始全面更新数据库注释...")

	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	sqls := []string{
		// 1. 系统用户表 sys_user
		"ALTER TABLE sys_user COMMENT = '用户信息表'",
		"ALTER TABLE sys_user MODIFY COLUMN user_id BIGINT AUTO_INCREMENT COMMENT '用户ID'",
		"ALTER TABLE sys_user MODIFY COLUMN dept_id BIGINT COMMENT '部门ID'",
		"ALTER TABLE sys_user MODIFY COLUMN user_name VARCHAR(30) NOT NULL COMMENT '用户账号'",
		"ALTER TABLE sys_user MODIFY COLUMN nick_name VARCHAR(30) COMMENT '用户昵称'",
		"ALTER TABLE sys_user MODIFY COLUMN password VARCHAR(100) NOT NULL COMMENT '密码'",
		"ALTER TABLE sys_user MODIFY COLUMN email VARCHAR(50) COMMENT '邮箱'",
		"ALTER TABLE sys_user MODIFY COLUMN phone VARCHAR(11) COMMENT '手机号'",
		"ALTER TABLE sys_user MODIFY COLUMN avatar VARCHAR(255) COMMENT '头像'",
		"ALTER TABLE sys_user MODIFY COLUMN status CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)'",
		"ALTER TABLE sys_user MODIFY COLUMN del_flag CHAR(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)'",
		"ALTER TABLE sys_user MODIFY COLUMN login_ip VARCHAR(128) COMMENT '最后登录IP'",
		"ALTER TABLE sys_user MODIFY COLUMN login_date DATETIME COMMENT '最后登录时间'",
		"ALTER TABLE sys_user MODIFY COLUMN create_by BIGINT COMMENT '创建者'",
		"ALTER TABLE sys_user MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_user MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		// 2. 角色表 sys_role
		"ALTER TABLE sys_role COMMENT = '角色信息表'",
		"ALTER TABLE sys_role MODIFY COLUMN role_id BIGINT AUTO_INCREMENT COMMENT '角色ID'",
		"ALTER TABLE sys_role MODIFY COLUMN role_name VARCHAR(30) NOT NULL COMMENT '角色名称'",
		"ALTER TABLE sys_role MODIFY COLUMN role_key VARCHAR(100) NOT NULL COMMENT '角色权限字符'",
		"ALTER TABLE sys_role MODIFY COLUMN data_scope CHAR(1) DEFAULT '1' COMMENT '数据范围(1全部 2自定义 3本部门及以下 4本部门 5仅本人)'",
		"ALTER TABLE sys_role MODIFY COLUMN sort INT DEFAULT 0 COMMENT '显示顺序'",
		"ALTER TABLE sys_role MODIFY COLUMN status CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)'",
		"ALTER TABLE sys_role MODIFY COLUMN del_flag CHAR(1) DEFAULT '0' COMMENT '删除标志'",
		"ALTER TABLE sys_role MODIFY COLUMN remark VARCHAR(500) COMMENT '备注'",
		"ALTER TABLE sys_role MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_role MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		// 3. 部门表 sys_dept
		"ALTER TABLE sys_dept COMMENT = '部门表'",
		"ALTER TABLE sys_dept MODIFY COLUMN dept_id BIGINT AUTO_INCREMENT COMMENT '部门ID'",
		"ALTER TABLE sys_dept MODIFY COLUMN parent_id BIGINT DEFAULT 0 COMMENT '父部门ID'",
		"ALTER TABLE sys_dept MODIFY COLUMN ancestors VARCHAR(500) COMMENT '祖级列表'",
		"ALTER TABLE sys_dept MODIFY COLUMN dept_name VARCHAR(30) NOT NULL COMMENT '部门名称'",
		"ALTER TABLE sys_dept MODIFY COLUMN sort INT DEFAULT 0 COMMENT '显示顺序'",
		"ALTER TABLE sys_dept MODIFY COLUMN leader VARCHAR(20) COMMENT '负责人'",
		"ALTER TABLE sys_dept MODIFY COLUMN phone VARCHAR(11) COMMENT '联系电话'",
		"ALTER TABLE sys_dept MODIFY COLUMN email VARCHAR(50) COMMENT '邮箱'",
		"ALTER TABLE sys_dept MODIFY COLUMN status CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)'",
		"ALTER TABLE sys_dept MODIFY COLUMN del_flag CHAR(1) DEFAULT '0' COMMENT '删除标志'",
		"ALTER TABLE sys_dept MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_dept MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		// 4. 菜单表 sys_menu
		"ALTER TABLE sys_menu COMMENT = '菜单权限表'",
		"ALTER TABLE sys_menu MODIFY COLUMN menu_id BIGINT AUTO_INCREMENT COMMENT '菜单ID'",
		"ALTER TABLE sys_menu MODIFY COLUMN parent_id BIGINT DEFAULT 0 COMMENT '父菜单ID'",
		"ALTER TABLE sys_menu MODIFY COLUMN menu_name VARCHAR(50) NOT NULL COMMENT '菜单名称'",
		"ALTER TABLE sys_menu MODIFY COLUMN menu_type CHAR(1) COMMENT '菜单类型(M目录 C菜单 F按钮)'",
		"ALTER TABLE sys_menu MODIFY COLUMN path VARCHAR(200) COMMENT '路由地址'",
		"ALTER TABLE sys_menu MODIFY COLUMN component VARCHAR(255) COMMENT '组件路径'",
		"ALTER TABLE sys_menu MODIFY COLUMN perms VARCHAR(100) COMMENT '权限标识'",
		"ALTER TABLE sys_menu MODIFY COLUMN icon VARCHAR(100) COMMENT '菜单图标'",
		"ALTER TABLE sys_menu MODIFY COLUMN sort INT DEFAULT 0 COMMENT '显示顺序'",
		"ALTER TABLE sys_menu MODIFY COLUMN visible CHAR(1) DEFAULT '0' COMMENT '显示状态(0显示 1隐藏)'",
		"ALTER TABLE sys_menu MODIFY COLUMN status CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)'",
		"ALTER TABLE sys_menu MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_menu MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		// 5. 用户角色关联表 sys_user_role
		"ALTER TABLE sys_user_role COMMENT = '用户角色关联表'",
		"ALTER TABLE sys_user_role MODIFY COLUMN user_id BIGINT NOT NULL COMMENT '用户ID'",
		"ALTER TABLE sys_user_role MODIFY COLUMN role_id BIGINT NOT NULL COMMENT '角色ID'",

		// 6. 角色菜单关联表 sys_role_menu
		"ALTER TABLE sys_role_menu COMMENT = '角色菜单关联表'",
		"ALTER TABLE sys_role_menu MODIFY COLUMN role_id BIGINT NOT NULL COMMENT '角色ID'",
		"ALTER TABLE sys_role_menu MODIFY COLUMN menu_id BIGINT NOT NULL COMMENT '菜单ID'",

		// 7. 角色部门关联表 sys_role_dept
		"ALTER TABLE sys_role_dept COMMENT = '角色部门关联表'",
		"ALTER TABLE sys_role_dept MODIFY COLUMN role_id BIGINT NOT NULL COMMENT '角色ID'",
		"ALTER TABLE sys_role_dept MODIFY COLUMN dept_id BIGINT NOT NULL COMMENT '部门ID'",

		// 8. 操作日志表 sys_oper_log
		"ALTER TABLE sys_oper_log COMMENT = '操作日志记录'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_id BIGINT AUTO_INCREMENT COMMENT '日志ID'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN title VARCHAR(50) COMMENT '模块标题'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN business_type INT DEFAULT 0 COMMENT '业务类型(0其他 1新增 2修改 3删除)'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN method VARCHAR(100) COMMENT '方法名称'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN request_method VARCHAR(10) COMMENT '请求方式'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_user_id BIGINT COMMENT '操作人员ID'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_user_name VARCHAR(50) COMMENT '操作人员'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_url VARCHAR(255) COMMENT '请求URL'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_ip VARCHAR(128) COMMENT '操作IP'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_param TEXT COMMENT '请求参数'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN json_result TEXT COMMENT '返回参数'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN status CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1异常)'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN error_msg TEXT COMMENT '错误消息'",
		"ALTER TABLE sys_oper_log MODIFY COLUMN oper_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间'",
	}

	for _, sql := range sqls {
		// 简单的错误处理，打印错误但继续执行，因为某些列可能没有变化
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("执行警告/错误: %s, 错误信息: %v\n", sql, err)
		} else {
			// fmt.Printf("成功执行: %s\n", sql)
		}
	}

	fmt.Println("所有表和字段注释更新完成!")
}
