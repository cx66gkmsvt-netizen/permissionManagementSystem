package main

import (
	"fmt"
	"log"
	"time"
	"user-center/internal/model"
	"user-center/internal/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DSN = "root:Zhixin123..@tcp(43.138.1.141:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	fmt.Println("=== 开始测试跟进记录功能 ===")

	// 1. 连接数据库
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 2. 注入数据库实例 & 自动迁移
	repository.DB = db
	// 始终执行 AutoMigrate 以确保字段更新
	fmt.Println("正在执行 AutoMigrate (sys_follow_up) 并添加表注释...")
	err = db.Set("gorm:table_options", "COMMENT='跟进记录表'").AutoMigrate(&model.SysFollowUp{})
	if err != nil {
		log.Printf("AutoMigrate error: %v", err)
	}
	// 强制更新表注释 (GORM AutoMigrate 有时在表已存在时不更新 Table Option)
	db.Exec("ALTER TABLE sys_follow_up COMMENT = '跟进记录表'")

	// 3. 初始化 Repo
	repo := repository.NewFollowUpRepository()

	// 4. 创建测试数据
	userID := int64(1001)
	record := &model.SysFollowUp{
		TargetType:   "sys_user",
		TargetID:     userID,
		Content:      "测试跟进记录功能 - 修改了用户资料",
		OperUserID:   new(int64),
		OperUserName: "admin",
		OperTime:     time.Now(),
		Remark:       "这是一个测试备注",
	}
	*record.OperUserID = 1

	fmt.Println("正在写入跟进记录...")
	err = repo.Create(record)
	if err != nil {
		log.Fatalf("写入失败: %v", err)
	}
	fmt.Printf("写入成功! ID: %d\n", record.ID)

	// 5. 查询验证
	fmt.Println("正在查询跟进记录...")
	list, err := repo.ListByTarget("sys_user", userID)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	found := false
	for _, item := range list {
		fmt.Printf("记录ID: %d, 内容: %s, 备注: %s, 操作人: %s\n", item.ID, item.Content, item.Remark, item.OperUserName)
		if item.ID == record.ID && item.Remark == "这是一个测试备注" {
			found = true
		}
	}

	if found {
		fmt.Println("SUCCESS: 成功查询到刚写入的记录!")
	} else {
		fmt.Println("ERROR: 未找到写入的记录!")
	}

	// 6. 验证表注释
	var tableComment string
	db.Raw("SELECT TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sys_follow_up'").Scan(&tableComment)
	fmt.Printf("当前表注释: %s\n", tableComment)

	fmt.Println("=== 测试结束 ===")
}
