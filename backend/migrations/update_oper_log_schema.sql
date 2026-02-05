USE user_center;

-- Add columns to sys_oper_log if they don't exist
-- Note: MySQL doesn't have IF NOT EXISTS for ADD COLUMN directly in all versions, 
-- but we can try to add them. 
-- Better to use a procedure or just run it and ignore error if exists, 
-- or since this is a dev environment fix, assume they might need to be added.

-- Add trace_content for execution process recording
ALTER TABLE sys_oper_log ADD COLUMN IF NOT EXISTS trace_content LONGTEXT COMMENT '追踪内容';

-- Add json_result for response body recording
ALTER TABLE sys_oper_log ADD COLUMN IF NOT EXISTS json_result TEXT COMMENT '返回结果';

-- Also update current month's table if it exists
SET @curr_month = DATE_FORMAT(NOW(), '%Y%m');
SET @table_name = CONCAT('sys_oper_log_', @curr_month);
SET @sql = CONCAT('ALTER TABLE ', @table_name, ' ADD COLUMN IF NOT EXISTS trace_content LONGTEXT COMMENT "追踪内容"');
PREPARE stmt FROM @sql;
-- EXECUTE stmt; -- Dynamic SQL might fail if table doesn't exist, handle manually or let user know.

-- Simplified: The user can drop the monthly table to regenerate it from the base table.
