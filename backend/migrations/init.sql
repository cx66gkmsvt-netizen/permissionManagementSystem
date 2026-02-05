-- =====================================================
-- 用户中心权限管理系统 - 数据库初始化脚本
-- MySQL 8.0+
-- =====================================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS user_center DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE user_center;

-- =====================================================
-- 1. 用户表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_user (
  user_id      BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
  dept_id      BIGINT COMMENT '部门ID',
  user_name    VARCHAR(30) NOT NULL COMMENT '用户账号',
  nick_name    VARCHAR(30) COMMENT '用户昵称',
  password     VARCHAR(100) NOT NULL COMMENT '密码',
  email        VARCHAR(50) COMMENT '邮箱',
  phone        VARCHAR(11) COMMENT '手机号',
  avatar       VARCHAR(255) COMMENT '头像',
  status       CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  del_flag     CHAR(1) DEFAULT '0' COMMENT '删除标志(0存在 2删除)',
  login_ip     VARCHAR(128) COMMENT '最后登录IP',
  login_date   DATETIME COMMENT '最后登录时间',
  create_by    BIGINT COMMENT '创建者',
  create_time  DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY idx_user_name (user_name)
) ENGINE=InnoDB COMMENT='用户信息表';

-- =====================================================
-- 2. 角色表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_role (
  role_id      BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
  role_name    VARCHAR(30) NOT NULL COMMENT '角色名称',
  role_key     VARCHAR(100) NOT NULL COMMENT '角色权限字符',
  data_scope   CHAR(1) DEFAULT '1' COMMENT '数据范围(1全部 2自定义 3本部门及以下 4本部门 5仅本人)',
  sort         INT DEFAULT 0 COMMENT '显示顺序',
  status       CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  del_flag     CHAR(1) DEFAULT '0' COMMENT '删除标志',
  remark       VARCHAR(500) COMMENT '备注',
  create_time  DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY idx_role_key (role_key)
) ENGINE=InnoDB COMMENT='角色信息表';

-- =====================================================
-- 3. 部门表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_dept (
  dept_id      BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '部门ID',
  parent_id    BIGINT DEFAULT 0 COMMENT '父部门ID',
  ancestors    VARCHAR(500) COMMENT '祖级列表',
  dept_name    VARCHAR(30) NOT NULL COMMENT '部门名称',
  sort         INT DEFAULT 0 COMMENT '显示顺序',
  leader       VARCHAR(20) COMMENT '负责人',
  phone        VARCHAR(11) COMMENT '联系电话',
  email        VARCHAR(50) COMMENT '邮箱',
  status       CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  del_flag     CHAR(1) DEFAULT '0' COMMENT '删除标志',
  create_time  DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB COMMENT='部门表';

-- =====================================================
-- 4. 菜单表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_menu (
  menu_id      BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '菜单ID',
  parent_id    BIGINT DEFAULT 0 COMMENT '父菜单ID',
  menu_name    VARCHAR(50) NOT NULL COMMENT '菜单名称',
  menu_type    CHAR(1) COMMENT '菜单类型(M目录 C菜单 F按钮)',
  path         VARCHAR(200) COMMENT '路由地址',
  component    VARCHAR(255) COMMENT '组件路径',
  perms        VARCHAR(100) COMMENT '权限标识',
  icon         VARCHAR(100) COMMENT '菜单图标',
  sort         INT DEFAULT 0 COMMENT '显示顺序',
  visible      CHAR(1) DEFAULT '0' COMMENT '显示状态(0显示 1隐藏)',
  status       CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1停用)',
  create_time  DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB COMMENT='菜单权限表';

-- =====================================================
-- 5. 用户-角色关联表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_user_role (
  user_id      BIGINT NOT NULL COMMENT '用户ID',
  role_id      BIGINT NOT NULL COMMENT '角色ID',
  PRIMARY KEY (user_id, role_id)
) ENGINE=InnoDB COMMENT='用户角色关联表';

-- =====================================================
-- 6. 角色-菜单关联表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_role_menu (
  role_id      BIGINT NOT NULL COMMENT '角色ID',
  menu_id      BIGINT NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (role_id, menu_id)
) ENGINE=InnoDB COMMENT='角色菜单关联表';

-- =====================================================
-- 7. 角色-部门关联表 (数据权限)
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_role_dept (
  role_id      BIGINT NOT NULL COMMENT '角色ID',
  dept_id      BIGINT NOT NULL COMMENT '部门ID',
  PRIMARY KEY (role_id, dept_id)
) ENGINE=InnoDB COMMENT='角色部门关联表';

-- =====================================================
-- 8. 操作日志表
-- =====================================================
CREATE TABLE IF NOT EXISTS sys_oper_log (
  oper_id        BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
  title          VARCHAR(50) COMMENT '模块标题',
  business_type  INT DEFAULT 0 COMMENT '业务类型(0其他 1新增 2修改 3删除)',
  method         VARCHAR(100) COMMENT '方法名称',
  request_method VARCHAR(10) COMMENT '请求方式',
  oper_user_id   BIGINT COMMENT '操作人员ID',
  oper_user_name VARCHAR(50) COMMENT '操作人员',
  oper_url       VARCHAR(255) COMMENT '请求URL',
  oper_ip        VARCHAR(128) COMMENT '操作IP',
  oper_param     TEXT COMMENT '请求参数',
  json_result    TEXT COMMENT '返回参数',
  status         CHAR(1) DEFAULT '0' COMMENT '状态(0正常 1异常)',
  error_msg      TEXT COMMENT '错误消息',
  oper_time      DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间'
) ENGINE=InnoDB COMMENT='操作日志记录';

-- =====================================================
-- 初始化数据
-- =====================================================

-- 默认部门
INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort, status) VALUES
(1, 0, '0', '总公司', 0, '0'),
(2, 1, '0,1', '研发部', 1, '0'),
(3, 1, '0,1', '市场部', 2, '0'),
(4, 1, '0,1', '财务部', 3, '0');

-- 默认角色
INSERT INTO sys_role (role_id, role_name, role_key, data_scope, sort, status) VALUES
(1, '超级管理员', 'admin', '1', 1, '0'),
(2, '普通角色', 'common', '5', 2, '0');

-- 默认用户 (密码: admin123)
-- BCrypt hash for 'admin123'
INSERT INTO sys_user (user_id, dept_id, user_name, nick_name, password, status) VALUES
(1, 1, 'admin', '超级管理员', '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '0');

-- 用户角色关联
INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1);

-- 默认菜单
INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) VALUES
-- 系统管理目录
(1, 0, '系统管理', 'M', '/system', '', '', 'Setting', 1, '0', '0'),
-- 用户管理
(2, 1, '用户管理', 'C', 'user', 'system/user/index', 'system:user:list', 'User', 1, '0', '0'),
(100, 2, '用户新增', 'F', '', '', 'system:user:add', '', 1, '0', '0'),
(101, 2, '用户修改', 'F', '', '', 'system:user:edit', '', 2, '0', '0'),
(102, 2, '用户删除', 'F', '', '', 'system:user:remove', '', 3, '0', '0'),
-- 角色管理
(3, 1, '角色管理', 'C', 'role', 'system/role/index', 'system:role:list', 'UserFilled', 2, '0', '0'),
(110, 3, '角色新增', 'F', '', '', 'system:role:add', '', 1, '0', '0'),
(111, 3, '角色修改', 'F', '', '', 'system:role:edit', '', 2, '0', '0'),
(112, 3, '角色删除', 'F', '', '', 'system:role:remove', '', 3, '0', '0'),
-- 菜单管理
(4, 1, '菜单管理', 'C', 'menu', 'system/menu/index', 'system:menu:list', 'Menu', 3, '0', '0'),
(120, 4, '菜单新增', 'F', '', '', 'system:menu:add', '', 1, '0', '0'),
(121, 4, '菜单修改', 'F', '', '', 'system:menu:edit', '', 2, '0', '0'),
(122, 4, '菜单删除', 'F', '', '', 'system:menu:remove', '', 3, '0', '0'),
-- 部门管理
(5, 1, '部门管理', 'C', 'dept', 'system/dept/index', 'system:dept:list', 'OfficeBuilding', 4, '0', '0'),
(130, 5, '部门新增', 'F', '', '', 'system:dept:add', '', 1, '0', '0'),
(131, 5, '部门修改', 'F', '', '', 'system:dept:edit', '', 2, '0', '0'),
(132, 5, '部门删除', 'F', '', '', 'system:dept:remove', '', 3, '0', '0');

-- 角色菜单关联 (超级管理员拥有所有菜单)
INSERT INTO sys_role_menu (role_id, menu_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
(1, 100), (1, 101), (1, 102),
(1, 110), (1, 111), (1, 112),
(1, 120), (1, 121), (1, 122),
(1, 130), (1, 131), (1, 132);

-- 普通角色菜单关联 (只读权限)
INSERT INTO sys_role_menu (role_id, menu_id) VALUES
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

-- =====================================================
-- CC管理菜单
-- =====================================================

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
