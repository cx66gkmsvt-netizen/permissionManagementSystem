-- 插入CC管理菜单
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) 
VALUES (6, 1, 'CC管理', 'C', 'cc', 'cc/index', 'system:cc:list', 'Service', 5, '0', '0')
ON DUPLICATE KEY UPDATE menu_name=menu_name;

-- 插入CC管理按钮权限
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort, visible, status) 
VALUES (140, 6, 'CC新增', 'F', 'system:cc:add', 1, '0', '0')
ON DUPLICATE KEY UPDATE menu_name=menu_name;

INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort, visible, status) 
VALUES (141, 6, 'CC修改', 'F', 'system:cc:edit', 2, '0', '0')
ON DUPLICATE KEY UPDATE menu_name=menu_name;

INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, perms, sort, visible, status) 
VALUES (142, 6, 'CC删除', 'F', 'system:cc:remove', 3, '0', '0')
ON DUPLICATE KEY UPDATE menu_name=menu_name;

-- 给管理员角色添加CC菜单权限
INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 6) ON DUPLICATE KEY UPDATE role_id=role_id;
INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 140) ON DUPLICATE KEY UPDATE role_id=role_id;
INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 141) ON DUPLICATE KEY UPDATE role_id=role_id;
INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, 142) ON DUPLICATE KEY UPDATE role_id=role_id;
