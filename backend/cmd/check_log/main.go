package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SysOperLog 操作日志表
type SysOperLog struct {
	OperID   int64     `json:"operId" gorm:"primaryKey;autoIncrement;column:oper_id"`
	Title    string    `json:"title" gorm:"column:title;size:50"`
	OperTime time.Time `json:"operTime" gorm:"column:oper_time;autoCreateTime"`
}

func main() {
	// DSN from config or main.go
	dsn := "root:Zhixin123..@tcp(43.138.1.141:3306)/user_center?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	tableName := fmt.Sprintf("sys_oper_log_%s", time.Now().Format("200601"))
	fmt.Printf("Checking table: %s\n", tableName)

	var count int64
	// Check if table exists
	if !db.Migrator().HasTable(tableName) {
		fmt.Printf("Table %s does not exist.\n", tableName)
		// Check base table
		if db.Migrator().HasTable("sys_oper_log") {
			fmt.Println("Base table sys_oper_log exists.")
		} else {
			fmt.Println("Base table sys_oper_log DOES NOT exist.")
		}
		return
	}

	db.Table(tableName).Count(&count)
	fmt.Printf("Total logs in %s: %d\n", tableName, count)

	var logs []SysOperLog
	db.Table(tableName).Order("oper_id DESC").Limit(5).Find(&logs)
	for _, l := range logs {
		fmt.Printf("ID: %d, Title: %s, Time: %s\n", l.OperID, l.Title, l.OperTime)
	}
}
