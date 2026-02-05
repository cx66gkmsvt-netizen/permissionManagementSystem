-- CC管理系统数据库迁移脚本
-- 版本: 1.0.0
-- 说明: 添加CC组织架构管理相关表

-- 1. 军团表
CREATE TABLE IF NOT EXISTS `cc_legion` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '军团ID',
  `legion_name` varchar(50) NOT NULL COMMENT '军团名称',
  `leader_id` bigint DEFAULT NULL COMMENT '军团长ID',
  `balance` bigint DEFAULT 0 COMMENT '军团余额(分)',
  `status` char(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  `del_flag` char(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)',
  `create_by` varchar(50) DEFAULT NULL COMMENT '创建者',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_legion_name` (`legion_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC军团表';

-- 2. 团队表
CREATE TABLE IF NOT EXISTS `cc_team` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '团队ID',
  `team_name` varchar(50) NOT NULL COMMENT '团队名称',
  `business_type` varchar(20) DEFAULT NULL COMMENT '业务类型',
  `leader_id` bigint DEFAULT NULL COMMENT '团长ID',
  `legion_id` bigint DEFAULT NULL COMMENT '所属军团ID',
  `balance` bigint DEFAULT 0 COMMENT '团队余额(分)',
  `status` char(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  `del_flag` char(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)',
  `create_by` varchar(50) DEFAULT NULL COMMENT '创建者',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_legion_id` (`legion_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC团队表';

-- 3. 战队表
CREATE TABLE IF NOT EXISTS `cc_squad` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '战队ID',
  `squad_name` varchar(50) NOT NULL COMMENT '战队名称',
  `leader_id` bigint DEFAULT NULL COMMENT '战队长ID',
  `team_id` bigint NOT NULL COMMENT '所属团队ID',
  `balance` bigint DEFAULT 0 COMMENT '战队余额(分)',
  `status` char(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  `del_flag` char(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)',
  `create_by` varchar(50) DEFAULT NULL COMMENT '创建者',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_squad_name` (`squad_name`),
  KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC战队表';

-- 4. CC管理日志表
CREATE TABLE IF NOT EXISTS `cc_manage_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `log_type` varchar(50) NOT NULL COMMENT '日志类型',
  `target_type` varchar(20) DEFAULT NULL COMMENT '目标类型(legion/team/squad/cc)',
  `target_id` bigint DEFAULT NULL COMMENT '目标ID',
  `content` text COMMENT '日志内容',
  `operator_id` bigint DEFAULT NULL COMMENT '操作人ID',
  `operator_name` varchar(50) DEFAULT NULL COMMENT '操作人名称',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_target` (`target_type`, `target_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC管理日志表';

-- 5. CC资金日志表
CREATE TABLE IF NOT EXISTS `cc_fund_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `log_type` varchar(50) NOT NULL COMMENT '日志类型',
  `target_type` varchar(20) NOT NULL COMMENT '目标类型(legion/team/squad/cc)',
  `target_id` bigint NOT NULL COMMENT '目标ID',
  `target_name` varchar(50) DEFAULT NULL COMMENT '目标名称',
  `amount` bigint DEFAULT 0 COMMENT '变动金额(分)',
  `balance_before` bigint DEFAULT 0 COMMENT '变动前余额(分)',
  `balance_after` bigint DEFAULT 0 COMMENT '变动后余额(分)',
  `related_cc_id` bigint DEFAULT NULL COMMENT '关联CC ID',
  `related_squad_id` bigint DEFAULT NULL COMMENT '关联战队ID',
  `related_team_id` bigint DEFAULT NULL COMMENT '关联团队ID',
  `related_legion_id` bigint DEFAULT NULL COMMENT '关联军团ID',
  `content` text COMMENT '日志内容详情',
  `reason` varchar(500) DEFAULT NULL COMMENT '变动原因',
  `operator_id` bigint DEFAULT NULL COMMENT '操作人ID',
  `operator_name` varchar(50) DEFAULT NULL COMMENT '操作人名称',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_target` (`target_type`, `target_id`),
  KEY `idx_log_type` (`log_type`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC资金日志表';

-- 6. CC在班记录表
CREATE TABLE IF NOT EXISTS `cc_attendance` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `cc_id` bigint NOT NULL COMMENT 'CC成员ID',
  `attendance_date` date NOT NULL COMMENT '在班日期',
  `status` char(1) NOT NULL COMMENT '在班状态(1在班 2休班 3请假)',
  `operator_id` bigint DEFAULT NULL COMMENT '操作人ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cc_date` (`cc_id`, `attendance_date`),
  KEY `idx_attendance_date` (`attendance_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC在班记录表';

-- 7. CC例子分配表
CREATE TABLE IF NOT EXISTS `cc_lead_allocation` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `cc_id` bigint NOT NULL COMMENT 'CC成员ID',
  `allocation_date` date NOT NULL COMMENT '分配日期',
  `expected_allocation` int DEFAULT 0 COMMENT '预计分配',
  `actual_allocation` int DEFAULT 0 COMMENT '实际分配',
  `expected_supplement` int DEFAULT 0 COMMENT '预计补发',
  `actual_supplement` int DEFAULT 0 COMMENT '实际补发',
  `overdraft` int DEFAULT 0 COMMENT '透支数量',
  `processed_overdraft` int DEFAULT 0 COMMENT '已处理透支',
  `pending_overdraft` int DEFAULT 0 COMMENT '待处理透支',
  `is_allocated` char(1) DEFAULT '0' COMMENT '是否分配例子(0否 1是)',
  `allocation_rule` varchar(10) DEFAULT NULL COMMENT '分配规则(A节假日 B工作日无补偿 C工作日)',
  `allocation_reason` text COMMENT '分配/未分配原因',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cc_date` (`cc_id`, `allocation_date`),
  KEY `idx_allocation_date` (`allocation_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC例子分配表';

-- 8. CC资金配置表
CREATE TABLE IF NOT EXISTS `cc_fund_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_type` varchar(50) NOT NULL COMMENT '配置类型',
  `rank_or_id` varchar(50) DEFAULT NULL COMMENT 'ID/排名序号',
  `amount` bigint DEFAULT 0 COMMENT '金额(分)',
  `config_month` varchar(7) DEFAULT NULL COMMENT '配置月份(YYYY-MM)',
  `create_by` varchar(50) DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_config_type` (`config_type`),
  KEY `idx_config_month` (`config_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CC资金配置表';

-- 9. 修改cc_member表结构，添加新字段
ALTER TABLE `cc_member`
  ADD COLUMN IF NOT EXISTS `role_type` varchar(20) DEFAULT 'cc' COMMENT '角色类型(cc/squad_leader/team_leader/legion_leader)' AFTER `wechat`,
  ADD COLUMN IF NOT EXISTS `cno1` varchar(20) DEFAULT NULL COMMENT '天润CNO1' AFTER `role_type`,
  ADD COLUMN IF NOT EXISTS `cno2` varchar(20) DEFAULT NULL COMMENT '天润CNO2' AFTER `cno1`,
  ADD COLUMN IF NOT EXISTS `cloud_account1` varchar(50) DEFAULT NULL COMMENT '云客账号1' AFTER `cno2`,
  ADD COLUMN IF NOT EXISTS `cloud_account2` varchar(50) DEFAULT NULL COMMENT '云客账号2' AFTER `cloud_account1`,
  ADD COLUMN IF NOT EXISTS `ronglian_seat` varchar(50) DEFAULT NULL COMMENT '容联座席号' AFTER `cloud_account2`,
  ADD COLUMN IF NOT EXISTS `diankong_seat` varchar(50) DEFAULT NULL COMMENT '点控云座席号' AFTER `ronglian_seat`,
  ADD COLUMN IF NOT EXISTS `heli_account` varchar(50) DEFAULT NULL COMMENT '合力亿捷账号' AFTER `diankong_seat`,
  ADD COLUMN IF NOT EXISTS `baichuan_seat` varchar(50) DEFAULT NULL COMMENT '百川智通座席号' AFTER `heli_account`,
  ADD COLUMN IF NOT EXISTS `diankong_outbound_pool` int DEFAULT 1 COMMENT '点控云公海智能外呼(0关闭 1开启)' AFTER `baichuan_seat`,
  ADD COLUMN IF NOT EXISTS `diankong_outbound_list` int DEFAULT 0 COMMENT '点控云客户列表智能外呼(0关闭 1开启)' AFTER `diankong_outbound_pool`,
  ADD COLUMN IF NOT EXISTS `baichuan_outbound_pool` int DEFAULT 0 COMMENT '百川公海智能外呼(0关闭 1开启)' AFTER `diankong_outbound_list`,
  ADD COLUMN IF NOT EXISTS `baichuan_outbound_list` int DEFAULT 0 COMMENT '百川客户列表智能外呼(0关闭 1开启)' AFTER `baichuan_outbound_pool`,
  ADD COLUMN IF NOT EXISTS `legion_id` bigint DEFAULT NULL COMMENT '所属军团ID' AFTER `team_id`,
  ADD COLUMN IF NOT EXISTS `monthly_performance` bigint DEFAULT 0 COMMENT '当月业绩(分)' AFTER `balance`,
  ADD COLUMN IF NOT EXISTS `performance_rank` int DEFAULT 0 COMMENT '业绩排名' AFTER `monthly_performance`,
  ADD COLUMN IF NOT EXISTS `attendance_status` char(1) DEFAULT '2' COMMENT '在班状态(1在班 2休班 3请假)' AFTER `performance_rank`,
  ADD COLUMN IF NOT EXISTS `call_duration_target` int DEFAULT 0 COMMENT '通时目标(秒)' AFTER `attendance_status`,
  ADD COLUMN IF NOT EXISTS `call_count_target` int DEFAULT 0 COMMENT '通次目标' AFTER `call_duration_target`,
  ADD COLUMN IF NOT EXISTS `is_blocked` char(1) DEFAULT '0' COMMENT '是否屏蔽(0否 1是)' AFTER `status`;

-- 添加索引
ALTER TABLE `cc_member`
  ADD INDEX IF NOT EXISTS `idx_legion_id` (`legion_id`),
  ADD INDEX IF NOT EXISTS `idx_role_type` (`role_type`),
  ADD INDEX IF NOT EXISTS `idx_attendance_status` (`attendance_status`);
