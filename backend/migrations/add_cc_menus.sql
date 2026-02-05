-- =====================================================
-- CC管理菜单初始化脚本
-- =====================================================

USE user_center;

-- CC管理目录
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) VALUES
(6, 0, 'CC管理', 'M', '/cc', '', '', 'Avatar', 2, '0', '0');

-- CC管理子菜单
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) VALUES
-- CC列表
(7, 6, 'CC列表', 'C', '/system/cc', 'cc/index', 'system:cc:list', 'User', 1, '0', '0'),
-- 军团管理
(8, 6, '军团管理', 'C', '/cc/legion', 'cc/legion/index', 'cc:legion:list', 'Flag', 2, '0', '0'),
-- 团队管理
(9, 6, '团队管理', 'C', '/cc/team', 'cc/team/index', 'cc:team:list', 'UserFilled', 3, '0', '0'),
-- 战队管理
(10, 6, '战队管理', 'C', '/cc/squad', 'cc/squad/index', 'cc:squad:list', 'Aim', 4, '0', '0'),
-- 在班管理
(11, 6, '在班管理', 'C', '/cc/attendance', 'cc/attendance/index', 'cc:attendance:list', 'Clock', 5, '0', '0'),
-- 例子分配
(12, 6, '例子分配', 'C', '/cc/lead-allocation', 'cc/leadAllocation/index', 'cc:lead-allocation:list', 'Share', 6, '0', '0');

-- CC管理按钮权限
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) VALUES
-- CC列表权限
(200, 7, 'CC新增', 'F', '', '', 'system:cc:add', '', 1, '0', '0'),
(201, 7, 'CC修改', 'F', '', '', 'system:cc:edit', '', 2, '0', '0'),
(202, 7, 'CC删除', 'F', '', '', 'system:cc:remove', '', 3, '0', '0'),
-- 军团管理权限
(210, 8, '军团新增', 'F', '', '', 'cc:legion:add', '', 1, '0', '0'),
(211, 8, '军团修改', 'F', '', '', 'cc:legion:edit', '', 2, '0', '0'),
(212, 8, '军团删除', 'F', '', '', 'cc:legion:remove', '', 3, '0', '0'),
-- 团队管理权限
(220, 9, '团队新增', 'F', '', '', 'cc:team:add', '', 1, '0', '0'),
(221, 9, '团队修改', 'F', '', '', 'cc:team:edit', '', 2, '0', '0'),
(222, 9, '团队删除', 'F', '', '', 'cc:team:remove', '', 3, '0', '0'),
-- 战队管理权限
(230, 10, '战队新增', 'F', '', '', 'cc:squad:add', '', 1, '0', '0'),
(231, 10, '战队修改', 'F', '', '', 'cc:squad:edit', '', 2, '0', '0'),
(232, 10, '战队删除', 'F', '', '', 'cc:squad:remove', '', 3, '0', '0'),
-- 在班管理权限
(240, 11, '在班记录新增', 'F', '', '', 'cc:attendance:add', '', 1, '0', '0'),
(241, 11, '在班记录修改', 'F', '', '', 'cc:attendance:edit', '', 2, '0', '0'),
(242, 11, '在班记录删除', 'F', '', '', 'cc:attendance:remove', '', 3, '0', '0'),
(243, 11, '在班记录导出', 'F', '', '', 'cc:attendance:export', '', 4, '0', '0'),
-- 例子分配权限
(250, 12, '例子分配修改', 'F', '', '', 'cc:lead-allocation:edit', '', 1, '0', '0'),
(251, 12, '例子分配查看详情', 'F', '', '', 'cc:lead-allocation:detail', '', 2, '0', '0');

-- 为超级管理员角色添加CC管理菜单权限
INSERT INTO sys_role_menu (role_id, menu_id) VALUES
(1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11), (1, 12),
(1, 200), (1, 201), (1, 202),
(1, 210), (1, 211), (1, 212),
(1, 220), (1, 221), (1, 222),
(1, 230), (1, 231), (1, 232),
(1, 240), (1, 241), (1, 242), (1, 243),
(1, 250), (1, 251);

-- 提交事务
COMMIT;
