package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"user-center/internal/config"
	"user-center/internal/handler"
	"user-center/internal/middleware"
	"user-center/internal/pkg"
	"user-center/internal/repository"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化JWT
	pkg.InitJWT(cfg.JWT.Secret)

	// 初始化数据库
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 初始化数据
	initData()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Cors())

	// 初始化Handler
	authHandler := handler.NewAuthHandler(cfg)
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	deptHandler := handler.NewDeptHandler()
	menuHandler := handler.NewMenuHandler()
	statsHandler := handler.NewStatsHandler()
	profileHandler := handler.NewProfileHandler()

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/auth/login", authHandler.Login)
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth())
	{
		// 认证相关
		auth.POST("/auth/logout", authHandler.Logout)
		auth.GET("/auth/info", authHandler.GetUserInfo)
		auth.GET("/auth/routes", authHandler.GetRoutes)

		// 系统统计
		auth.GET("/system/stats", statsHandler.GetStats)

		// 个人中心
		auth.GET("/system/profile", profileHandler.GetProfile)
		auth.PUT("/system/profile", profileHandler.UpdateProfile)
		auth.PUT("/system/profile/password", profileHandler.UpdatePassword)

		// 用户管理
		userGroup := auth.Group("/system/user")
		{
			userGroup.GET("", userHandler.List)
			userGroup.GET("/:id", userHandler.Get)
			userGroup.POST("", middleware.OperLog("用户管理", 1), userHandler.Create)
			userGroup.PUT("/:id", middleware.OperLog("用户管理", 2), userHandler.Update)
			userGroup.DELETE("/:id", middleware.OperLog("用户管理", 3), userHandler.Delete)
			userGroup.PUT("/:id/resetPwd", middleware.OperLog("用户管理", 2), userHandler.ResetPassword)
		}

		// 角色管理
		roleGroup := auth.Group("/system/role")
		{
			roleGroup.GET("", roleHandler.List)
			roleGroup.GET("/all", roleHandler.SelectAll)
			roleGroup.GET("/:id", roleHandler.Get)
			roleGroup.POST("", middleware.OperLog("角色管理", 1), roleHandler.Create)
			roleGroup.PUT("/:id", middleware.OperLog("角色管理", 2), roleHandler.Update)
			roleGroup.DELETE("/:id", middleware.OperLog("角色管理", 3), roleHandler.Delete)
		}

		// 部门管理
		deptGroup := auth.Group("/system/dept")
		{
			deptGroup.GET("", deptHandler.List)
			deptGroup.GET("/all", deptHandler.SelectAll)
			deptGroup.GET("/:id", deptHandler.Get)
			deptGroup.POST("", middleware.OperLog("部门管理", 1), deptHandler.Create)
			deptGroup.PUT("/:id", middleware.OperLog("部门管理", 2), deptHandler.Update)
			deptGroup.DELETE("/:id", middleware.OperLog("部门管理", 3), deptHandler.Delete)
		}

		// 菜单管理
		menuGroup := auth.Group("/system/menu")
		{
			menuGroup.GET("", menuHandler.List)
			menuGroup.GET("/all", menuHandler.SelectAll)
			menuGroup.GET("/:id", menuHandler.Get)
			menuGroup.POST("", middleware.OperLog("菜单管理", 1), menuHandler.Create)
			menuGroup.PUT("/:id", middleware.OperLog("菜单管理", 2), menuHandler.Update)
			menuGroup.DELETE("/:id", middleware.OperLog("菜单管理", 3), menuHandler.Delete)
		}

		// CC管理
		ccHandler := handler.NewCCHandler()
		ccGroup := auth.Group("/system/cc")
		{
			ccGroup.GET("", ccHandler.List)
			ccGroup.GET("/:id", ccHandler.Get)
			ccGroup.POST("", middleware.OperLog("CC管理", 1), ccHandler.Create)
			ccGroup.PUT("/:id", middleware.OperLog("CC管理", 2), ccHandler.Update)
			ccGroup.DELETE("/:id", middleware.OperLog("CC管理", 3), ccHandler.Delete)
		}
	}

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initData 初始化默认数据
func initData() {
	db := repository.GetDB()

	// 检查是否已有管理员
	var count int64
	db.Table("sys_user").Count(&count)
	if count > 0 {
		return
	}

	// 创建默认管理员密码: admin123
	hashedPwd, _ := pkg.HashPassword("admin123")
	db.Exec(`INSERT INTO sys_user (user_id, user_name, nick_name, password, status) 
		VALUES (1, 'admin', '超级管理员', ?, '0')`, hashedPwd)

	// 创建默认角色
	db.Exec(`INSERT INTO sys_role (role_id, role_name, role_key, data_scope, status) 
		VALUES (1, '超级管理员', 'admin', '1', '0')`)

	// 关联用户和角色
	db.Exec(`INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1)`)

	// 创建默认部门
	db.Exec(`INSERT INTO sys_dept (dept_id, parent_id, ancestors, dept_name, sort, status) 
		VALUES (1, 0, '0', '总公司', 0, '0')`)

	// 创建默认菜单
	menus := []struct {
		id        int64
		parentID  int64
		name      string
		menuType  string
		path      string
		component string
		perms     string
		icon      string
		sort      int
	}{
		{1, 0, "系统管理", "M", "/system", "", "", "Setting", 1},
		{2, 1, "用户管理", "C", "user", "system/user/index", "system:user:list", "User", 1},
		{3, 1, "角色管理", "C", "role", "system/role/index", "system:role:list", "UserFilled", 2},
		{4, 1, "菜单管理", "C", "menu", "system/menu/index", "system:menu:list", "Menu", 3},
		{5, 1, "部门管理", "C", "dept", "system/dept/index", "system:dept:list", "OfficeBuilding", 4},
		// 用户管理按钮
		{100, 2, "用户新增", "F", "", "", "system:user:add", "", 1},
		{101, 2, "用户修改", "F", "", "", "system:user:edit", "", 2},
		{102, 2, "用户删除", "F", "", "", "system:user:remove", "", 3},
		// 角色管理按钮
		{110, 3, "角色新增", "F", "", "", "system:role:add", "", 1},
		{111, 3, "角色修改", "F", "", "", "system:role:edit", "", 2},
		{112, 3, "角色删除", "F", "", "", "system:role:remove", "", 3},
		// 菜单管理按钮
		{120, 4, "菜单新增", "F", "", "", "system:menu:add", "", 1},
		{121, 4, "菜单修改", "F", "", "", "system:menu:edit", "", 2},
		{122, 4, "菜单删除", "F", "", "", "system:menu:remove", "", 3},
		// 部门管理按钮
		{130, 5, "部门新增", "F", "", "", "system:dept:add", "", 1},
		{131, 5, "部门修改", "F", "", "", "system:dept:edit", "", 2},
		{132, 5, "部门删除", "F", "", "", "system:dept:remove", "", 3},
	}

	for _, m := range menus {
		db.Exec(`INSERT INTO sys_menu (menu_id, parent_id, menu_name, menu_type, path, component, perms, icon, sort, visible, status) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, '0', '0')`,
			m.id, m.parentID, m.name, m.menuType, m.path, m.component, m.perms, m.icon, m.sort)
		// 关联超管角色和菜单
		db.Exec(`INSERT INTO sys_role_menu (role_id, menu_id) VALUES (1, ?)`, m.id)
	}

	log.Println("Default data initialized successfully")
}
