package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Legion struct {
	ID   int
	Name string
}

type Team struct {
	ID       int
	Name     string
	LegionID int
}

type Squad struct {
	ID     int
	Name   string
	TeamID int
}

type CCMember struct {
	ID       int
	Name     string
	Mobile   string
	SquadID  int
	TeamID   int
	LegionID int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	f, err := os.Create("backend/migrations/seed_cc_data.sql")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	write := func(s string) {
		f.WriteString(s + "\n")
	}

	write("-- CC管理模块测试数据 (50条/表)")
	write("USE user_center;")
	write("SET NAMES utf8mb4;")

	// 1. Create cc_member table if not exists (Inferred from model)
	write(`
CREATE TABLE IF NOT EXISTS cc_member (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'CC成员ID',
  name varchar(50) NOT NULL COMMENT '老师名/姓名',
  nick_name varchar(50) DEFAULT NULL COMMENT '昵称',
  mobile varchar(11) NOT NULL COMMENT '手机号/账号',
  wechat varchar(50) DEFAULT NULL COMMENT '微信号',
  role_type varchar(20) DEFAULT 'cc' COMMENT '角色类型(cc/squad_leader/team_leader/legion_leader)',
  cno1 varchar(20) DEFAULT NULL COMMENT '天润CNO1',
  cno2 varchar(20) DEFAULT NULL COMMENT '天润CNO2',
  cloud_account1 varchar(50) DEFAULT NULL COMMENT '云客账号1',
  cloud_account2 varchar(50) DEFAULT NULL COMMENT '云客账号2',
  ronglian_seat varchar(50) DEFAULT NULL COMMENT '容联座席号',
  diankong_seat varchar(50) DEFAULT NULL COMMENT '点控云座席号',
  heli_account varchar(50) DEFAULT NULL COMMENT '合力亿捷账号',
  baichuan_seat varchar(50) DEFAULT NULL COMMENT '百川智通座席号',
  diankong_outbound_pool int DEFAULT 1 COMMENT '点控云公海智能外呼(0关闭 1开启)',
  diankong_outbound_list int DEFAULT 0 COMMENT '点控云客户列表智能外呼(0关闭 1开启)',
  baichuan_outbound_pool int DEFAULT 0 COMMENT '百川公海智能外呼(0关闭 1开启)',
  baichuan_outbound_list int DEFAULT 0 COMMENT '百川客户列表智能外呼(0关闭 1开启)',
  squad_id bigint DEFAULT NULL COMMENT '所属战队ID',
  team_id bigint DEFAULT NULL COMMENT '所属团队ID',
  legion_id bigint DEFAULT NULL COMMENT '所属军团ID',
  balance bigint DEFAULT 0 COMMENT '个人余额(分)',
  monthly_performance bigint DEFAULT 0 COMMENT '当月业绩(分)',
  performance_rank int DEFAULT 0 COMMENT '业绩排名',
  attendance_status char(1) DEFAULT '2' COMMENT '在班状态(1在班 2休班 3请假)',
  call_duration_target int DEFAULT 0 COMMENT '通时目标(秒)',
  call_count_target int DEFAULT 0 COMMENT '通次目标',
  status char(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  is_blocked char(1) DEFAULT '0' COMMENT '是否屏蔽(0否 1是)',
  del_flag char(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)',
  create_by varchar(50) DEFAULT NULL COMMENT '创建者',
  update_by varchar(50) DEFAULT NULL COMMENT '更新者',
  create_time datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_mobile (mobile),
  KEY idx_squad_id (squad_id),
  KEY idx_team_id (team_id),
  KEY idx_legion_id (legion_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC成员表';`)

	// Generate Data in Memory
	var legions []Legion
	for i := 1; i <= 50; i++ {
		legions = append(legions, Legion{ID: i, Name: fmt.Sprintf("军团%d", i)})
	}

	var teams []Team
	for i := 1; i <= 50; i++ {
		l := legions[rand.Intn(len(legions))]
		teams = append(teams, Team{ID: i, Name: fmt.Sprintf("团队%d", i), LegionID: l.ID})
	}

	var squads []Squad
	for i := 1; i <= 50; i++ {
		t := teams[rand.Intn(len(teams))]
		squads = append(squads, Squad{ID: i, Name: fmt.Sprintf("战队%d", i), TeamID: t.ID})
	}

	var members []CCMember
	for i := 1; i <= 50; i++ {
		s := squads[rand.Intn(len(squads))]
		// Find team for this squad
		var t Team
		for _, team := range teams {
			if team.ID == s.TeamID {
				t = team
				break
			}
		}
		// Find legion for this team
		var l Legion
		for _, legion := range legions {
			if legion.ID == t.LegionID {
				l = legion
				break
			}
		}

		members = append(members, CCMember{
			ID:       i,
			Name:     fmt.Sprintf("CC_%d", i),
			Mobile:   fmt.Sprintf("1380000%04d", i),
			SquadID:  s.ID,
			TeamID:   t.ID,
			LegionID: l.ID,
		})
	}

	// Write SQL
	write("\n-- Insert Legions")
	write("INSERT INTO cc_legion (id, legion_name, status) VALUES")
	for i, v := range legions {
		end := ","
		if i == len(legions)-1 {
			end = ";"
		}
		write(fmt.Sprintf("(%d, '%s', '0')%s", v.ID, v.Name, end))
	}

	write("\n-- Insert Teams")
	write("INSERT INTO cc_team (id, team_name, legion_id, status) VALUES")
	for i, v := range teams {
		end := ","
		if i == len(teams)-1 {
			end = ";"
		}
		write(fmt.Sprintf("(%d, '%s', %d, '0')%s", v.ID, v.Name, v.LegionID, end))
	}

	write("\n-- Insert Squads")
	write("INSERT INTO cc_squad (id, squad_name, team_id, status) VALUES")
	for i, v := range squads {
		end := ","
		if i == len(squads)-1 {
			end = ";"
		}
		write(fmt.Sprintf("(%d, '%s', %d, '0')%s", v.ID, v.Name, v.TeamID, end))
	}

	write("\n-- Insert CCMembers")
	write("INSERT INTO cc_member (id, name, mobile, squad_id, team_id, legion_id, status) VALUES")
	for i, v := range members {
		end := ","
		if i == len(members)-1 {
			end = ";"
		}
		write(fmt.Sprintf("(%d, '%s', '%s', %d, %d, %d, '0')%s", v.ID, v.Name, v.Mobile, v.SquadID, v.TeamID, v.LegionID, end))
	}

	write("\n-- Insert Attendance (for these members, today)")
	write("INSERT INTO cc_attendance (cc_id, attendance_date, status) VALUES")
	for i, v := range members {
		end := ","
		if i == len(members)-1 {
			end = ";"
		}
		status := rand.Intn(3) + 1 // 1, 2, 3
		write(fmt.Sprintf("(%d, CURRENT_DATE, '%d')%s", v.ID, status, end))
	}

	write("\n-- Insert Lead Allocation (for these members, today)")
	write("INSERT INTO cc_lead_allocation (cc_id, allocation_date, expected_allocation, actual_allocation, is_allocated) VALUES")
	for i, v := range members {
		end := ","
		if i == len(members)-1 {
			end = ";"
		}
		isAllocated := rand.Intn(2) // 0, 1
		write(fmt.Sprintf("(%d, CURRENT_DATE, 10, %d, '%d')%s", v.ID, rand.Intn(10), isAllocated, end))
	}

	fmt.Println("Seed SQL generated successfully!")
}
