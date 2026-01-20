# 用户中心与权限管理系统

企业级统一用户中心与 RBAC 权限管理系统，支持数据权限隔离。

## 技术栈

- **后端**: Go 1.21+ (Gin + GORM + JWT)
- **前端**: Vue 3 + Vite + Element Plus + Pinia
- **数据库**: MySQL 8.0+

## 项目结构

```
test03/
├── backend/                    # 后端 Go 项目
│   ├── cmd/server/main.go      # 入口文件
│   ├── internal/
│   │   ├── config/             # 配置
│   │   ├── model/              # 数据模型
│   │   ├── repository/         # 数据访问层
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # 中间件
│   │   └── pkg/                # 工具包
│   ├── migrations/             # 数据库迁移
│   └── go.mod
│
└── frontend/                   # 前端 Vue 项目
    ├── src/
    │   ├── api/                # API 请求
    │   ├── components/         # 通用组件
    │   ├── directives/         # 指令
    │   ├── router/             # 路由
    │   ├── stores/             # Pinia 状态
    │   ├── styles/             # 样式
    │   ├── utils/              # 工具函数
    │   └── views/              # 页面
    ├── package.json
    └── vite.config.js
```

## 快速开始

### 1. 数据库初始化

```bash
# 创建数据库并导入初始数据
mysql -u root -p < backend/migrations/init.sql
```

### 2. 启动后端

```bash
cd backend

# 安装依赖
go mod tidy

# 配置环境变量 (可选，使用默认值)
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=user_center

# 启动服务
go run cmd/server/main.go
```

后端服务将在 http://localhost:8080 启动。

### 3. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 http://localhost:3000 启动。

## 默认账号

| 账号 | 密码 | 角色 |
|------|------|------|
| admin | admin123 | 超级管理员 |

## 核心功能

### 认证模块
- JWT Token 认证
- 登录/登出
- 密码 BCrypt 加密

### RBAC 权限
- 用户管理 (CRUD)
- 角色管理 (配置菜单权限)
- 菜单管理 (目录/菜单/按钮三级结构)
- 部门管理 (树形结构)

### 数据权限
| 类型 | 说明 |
|------|------|
| 全部数据 | 无数据过滤 |
| 自定义数据 | 指定部门数据 |
| 本部门及以下 | 递归子部门数据 |
| 仅本部门 | 本部门数据 |
| 仅本人 | 个人创建的数据 |

### 操作审计
- 自动记录增删改操作
- 记录操作人、IP、时间、参数

## API 文档

### 认证
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `GET /api/auth/info` - 用户信息
- `GET /api/auth/routes` - 路由菜单

### 用户管理
- `GET /api/system/user` - 列表
- `POST /api/system/user` - 创建
- `PUT /api/system/user/:id` - 更新
- `DELETE /api/system/user/:id` - 删除

### 角色管理
- `GET /api/system/role` - 列表
- `POST /api/system/role` - 创建
- `PUT /api/system/role/:id` - 更新
- `DELETE /api/system/role/:id` - 删除

### 部门管理
- `GET /api/system/dept` - 部门树
- `POST /api/system/dept` - 创建
- `PUT /api/system/dept/:id` - 更新
- `DELETE /api/system/dept/:id` - 删除

### 菜单管理
- `GET /api/system/menu` - 菜单树
- `POST /api/system/menu` - 创建
- `PUT /api/system/menu/:id` - 更新
- `DELETE /api/system/menu/:id` - 删除

## 前端权限控制

使用 `v-permission` 指令控制按钮显示：

```vue
<el-button v-permission="'system:user:add'">新增</el-button>
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 8080 | 服务端口 |
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USER | root | 数据库用户 |
| DB_PASSWORD | 123456 | 数据库密码 |
| DB_NAME | user_center | 数据库名称 |
| JWT_SECRET | user-center-secret-key-2024 | JWT 密钥 |

## License

MIT
