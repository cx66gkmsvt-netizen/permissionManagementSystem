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
	fmt.Println("开始更新数据库注释...")

	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	sqls := []string{
		"ALTER TABLE sys_role COMMENT = '角色信息表'",
		"ALTER TABLE sys_role MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_role MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		"ALTER TABLE sys_dept COMMENT = '部门表'",
		"ALTER TABLE sys_dept MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_dept MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",

		"ALTER TABLE sys_menu COMMENT = '菜单权限表'",
		"ALTER TABLE sys_menu MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'",
		"ALTER TABLE sys_menu MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'",
	}

	for _, sql := range sqls {
		fmt.Printf("执行: %s\n", sql)
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("执行失败: %v\n", err)
		}
	}

	fmt.Println("数据库注释更新完成!")
}
