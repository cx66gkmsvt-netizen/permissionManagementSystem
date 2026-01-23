package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 先连接MySQL服务器创建数据库
	dsnWithoutDB := "root:Zhixin123..@tcp(43.138.1.141:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接MySQL失败: %v", err)
	}

	// 创建数据库
	db.Exec("CREATE DATABASE IF NOT EXISTS user_center DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	log.Println("✓ 数据库 user_center 创建成功")

	// 连接到user_center数据库
	dsn := "root:Zhixin123..@tcp(43.138.1.141:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 创建表
	createTables(db)
	log.Println("✓ 数据表创建成功")

	// 插入初始数据
	insertData(db)
	log.Println("✓ 初始数据插入成功")

	fmt.Println("\n========================================")
	fmt.Println("数据库初始化完成!")
	fmt.Println("默认账号: admin")
	fmt.Println("默认密码: admin123")
	fmt.Println("========================================")
}

func createTables(db *gorm.DB) {
	// 用户表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_user (
		user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		dept_id BIGINT,
		user_name VARCHAR(30) NOT NULL UNIQUE,
		nick_name VARCHAR(30),
		password VARCHAR(100) NOT NULL,
		email VARCHAR(50),
		phone VARCHAR(11),
		avatar VARCHAR(255),
		status CHAR(1) DEFAULT '0',
		del_flag CHAR(1) DEFAULT '0',
		login_ip VARCHAR(128),
		login_date DATETIME,
		create_by BIGINT,
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// 角色表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_role (
		role_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		role_name VARCHAR(30) NOT NULL,
		role_key VARCHAR(100) NOT NULL UNIQUE,
		data_scope CHAR(1) DEFAULT '1',
		sort INT DEFAULT 0,
		status CHAR(1) DEFAULT '0',
		del_flag CHAR(1) DEFAULT '0',
		remark VARCHAR(500),
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// 部门表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_dept (
		dept_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		parent_id BIGINT DEFAULT 0,
		ancestors VARCHAR(500),
		dept_name VARCHAR(30) NOT NULL,
		sort INT DEFAULT 0,
		leader VARCHAR(20),
		phone VARCHAR(11),
		email VARCHAR(50),
		status CHAR(1) DEFAULT '0',
		del_flag CHAR(1) DEFAULT '0',
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// 菜单表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_menu (
		menu_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		parent_id BIGINT DEFAULT 0,
		menu_name VARCHAR(50) NOT NULL,
		menu_type CHAR(1),
		path VARCHAR(200),
		component VARCHAR(255),
		perms VARCHAR(100),
		icon VARCHAR(100),
		sort INT DEFAULT 0,
		visible CHAR(1) DEFAULT '0',
		status CHAR(1) DEFAULT '0',
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// 关联表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_user_role (
		user_id BIGINT NOT NULL,
		role_id BIGINT NOT NULL,
		PRIMARY KEY (user_id, role_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	db.Exec(`CREATE TABLE IF NOT EXISTS sys_role_menu (
		role_id BIGINT NOT NULL,
		menu_id BIGINT NOT NULL,
		PRIMARY KEY (role_id, menu_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	db.Exec(`CREATE TABLE IF NOT EXISTS sys_role_dept (
		role_id BIGINT NOT NULL,
		dept_id BIGINT NOT NULL,
		PRIMARY KEY (role_id, dept_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// 操作日志表
	db.Exec(`CREATE TABLE IF NOT EXISTS sys_oper_log (
		oper_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(50),
		business_type INT DEFAULT 0,
		method VARCHAR(100),
		request_method VARCHAR(10),
		oper_user_id BIGINT,
		oper_user_name VARCHAR(50),
		oper_url VARCHAR(255),
		oper_ip VARCHAR(128),
		oper_param TEXT,
		json_result TEXT,
		status CHAR(1) DEFAULT '0',
		error_msg TEXT,
		oper_time DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// CC成员表
	db.Exec(`CREATE TABLE IF NOT EXISTS cc_member (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(50) NOT NULL,
		nick_name VARCHAR(50),
		mobile VARCHAR(11) NOT NULL,
		wechat VARCHAR(50),
		cno VARCHAR(20),
		cloud_account VARCHAR(50),
		team_id BIGINT,
		squad_id BIGINT,
		balance DECIMAL(10,2) DEFAULT 0,
		status CHAR(1) DEFAULT '0',
		del_flag CHAR(1) DEFAULT '0',
		create_by VARCHAR(50),
		update_by VARCHAR(50),
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
}

func insertData(db *gorm.DB) {
	// 检查是否已有数据
	var count int64
	db.Raw("SELECT COUNT(*) FROM sys_user").Scan(&count)
	if count > 0 {
		log.Println("数据已存在，跳过插入")
		return
	}

	// 密码加密
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// 插入部门
	db.Exec("INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort) VALUES (1, 0, '0', '总公司', 0)")
	db.Exec("INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort) VALUES (2, 1, '0,1', '研发部', 1)")
	db.Exec("INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort) VALUES (3, 1, '0,1', '市场部', 2)")
	db.Exec("INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort) VALUES (4, 1, '0,1', '财务部', 3)")

	// 插入角色
	db.Exec("INSERT INTO sys_role (role_id, role_name, role_key, data_scope, sort) VALUES (1, '超级管理员', 'admin', '1', 1)")
	db.Exec("INSERT INTO sys_role (role_id, role_name, role_key, data_scope, sort) VALUES (2, '普通角色', 'common', '5', 2)")

	// 插入用户
	db.Exec("INSERT INTO sys_user (user_id, dept_id, user_name, nick_name, password, status) VALUES (1, 1, 'admin', '超级管理员', ?, '0')", string(hashedPwd))

	// 用户角色关联
	db.Exec("INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1)")

	// 插入菜单
	menus := []string{
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (1, 0, '系统管理', 'M', '/system', '', '', 'Setting', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (2, 1, '用户管理', 'C', 'user', 'system/user/index', 'system:user:list', 'User', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (3, 1, '角色管理', 'C', 'role', 'system/role/index', 'system:role:list', 'UserFilled', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (4, 1, '菜单管理', 'C', 'menu', 'system/menu/index', 'system:menu:list', 'Menu', 3)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (5, 1, '部门管理', 'C', 'dept', 'system/dept/index', 'system:dept:list', 'OfficeBuilding', 4)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (100, 2, '用户新增', 'F', 'system:user:add', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (101, 2, '用户修改', 'F', 'system:user:edit', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (102, 2, '用户删除', 'F', 'system:user:remove', 3)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (110, 3, '角色新增', 'F', 'system:role:add', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (111, 3, '角色修改', 'F', 'system:role:edit', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (112, 3, '角色删除', 'F', 'system:role:remove', 3)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (120, 4, '菜单新增', 'F', 'system:menu:add', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (121, 4, '菜单修改', 'F', 'system:menu:edit', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (122, 4, '菜单删除', 'F', 'system:menu:remove', 3)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (130, 5, '部门新增', 'F', 'system:dept:add', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (131, 5, '部门修改', 'F', 'system:dept:edit', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (132, 5, '部门删除', 'F', 'system:dept:remove', 3)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort) VALUES (6, 1, 'CC管理', 'C', 'cc', 'cc/index', 'system:cc:list', 'Service', 5)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (140, 6, 'CC新增', 'F', 'system:cc:add', 1)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (141, 6, 'CC修改', 'F', 'system:cc:edit', 2)",
		"INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort) VALUES (142, 6, 'CC删除', 'F', 'system:cc:remove', 3)",
	}
	for _, sql := range menus {
		db.Exec(sql)
	}

	// 角色菜单关联
	roleMenus := []string{
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 1)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 2)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 3)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 4)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 5)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 100)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 101)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 102)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 110)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 111)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 112)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 120)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 121)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 122)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 130)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 131)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 132)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 6)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 140)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 141)",
		"INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 142)",
	}
	for _, sql := range roleMenus {
		db.Exec(sql)
	}
}
