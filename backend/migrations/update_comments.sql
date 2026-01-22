-- =====================================================
-- 补全数据库注释迁移脚本
-- =====================================================

USE user_center;

-- 1. 角色表 (sys_role)
ALTER TABLE sys_role COMMENT = '角色信息表';
ALTER TABLE sys_role MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE sys_role MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

-- 2. 部门表 (sys_dept)
ALTER TABLE sys_dept COMMENT = '部门表';
ALTER TABLE sys_dept MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE sys_dept MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

-- 3. 菜单表 (sys_menu)
ALTER TABLE sys_menu COMMENT = '菜单权限表';
ALTER TABLE sys_menu MODIFY COLUMN create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间';
ALTER TABLE sys_menu MODIFY COLUMN update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';
