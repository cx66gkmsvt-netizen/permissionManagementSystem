-- CC角色菜单权限配置
-- 为CC相关角色只配置CC管理菜单，不配置系统管理菜单

-- 1. 确保CC角色存在（如果不存在则创建）
INSERT INTO sys_role (role_id, role_name, role_key, sort, status, del_flag)
VALUES (100, 'CC', 'cc', 100, '0', '0')
ON DUPLICATE KEY UPDATE role_name=VALUES(role_name);

INSERT INTO sys_role (role_id, role_name, role_key, sort, status, del_flag)
VALUES (101, 'CC战队长', 'cc_squad_leader', 101, '0', '0')
ON DUPLICATE KEY UPDATE role_name=VALUES(role_name);

INSERT INTO sys_role (role_id, role_name, role_key, sort, status, del_flag)
VALUES (102, 'CC团长', 'cc_team_leader', 102, '0', '0')
ON DUPLICATE KEY UPDATE role_name=VALUES(role_name);

INSERT INTO sys_role (role_id, role_name, role_key, sort, status, del_flag)
VALUES (103, 'CC军团长', 'cc_legion_leader', 103, '0', '0')
ON DUPLICATE KEY UPDATE role_name=VALUES(role_name);

-- 2. 清除CC角色原有菜单关联
DELETE FROM sys_role_menu WHERE role_id IN (100, 101, 102, 103);

-- 3. 为CC角色配置CC管理相关菜单
-- 注意：请根据实际数据库中的 menu_id 调整以下ID值
-- 可通过 SELECT menu_id, menu_name, path, perms FROM sys_menu ORDER BY menu_id; 查看

-- 将CC管理相关菜单关联到所有CC角色
-- 这里使用子查询动态获取菜单ID，避免硬编码

-- CC管理目录及其子菜单（路径包含 cc 相关的菜单）
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 100, menu_id FROM sys_menu WHERE perms LIKE 'cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 101, menu_id FROM sys_menu WHERE perms LIKE 'cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 102, menu_id FROM sys_menu WHERE perms LIKE 'cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 103, menu_id FROM sys_menu WHERE perms LIKE 'cc:%' AND status = '0';

-- 也需要包含system:cc相关权限的菜单（CC列表）
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 100, menu_id FROM sys_menu WHERE perms LIKE 'system:cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 101, menu_id FROM sys_menu WHERE perms LIKE 'system:cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 102, menu_id FROM sys_menu WHERE perms LIKE 'system:cc:%' AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 103, menu_id FROM sys_menu WHERE perms LIKE 'system:cc:%' AND status = '0';

-- 还需要加上父级目录菜单（menu_type='M'），这样侧边栏才能正确显示
-- CC管理目录
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 100, menu_id FROM sys_menu WHERE menu_type = 'M' AND (path LIKE '%cc%' OR menu_name LIKE '%CC%') AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 101, menu_id FROM sys_menu WHERE menu_type = 'M' AND (path LIKE '%cc%' OR menu_name LIKE '%CC%') AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 102, menu_id FROM sys_menu WHERE menu_type = 'M' AND (path LIKE '%cc%' OR menu_name LIKE '%CC%') AND status = '0';

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 103, menu_id FROM sys_menu WHERE menu_type = 'M' AND (path LIKE '%cc%' OR menu_name LIKE '%CC%') AND status = '0';
