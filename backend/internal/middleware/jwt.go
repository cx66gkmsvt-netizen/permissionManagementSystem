package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"user-center/internal/pkg"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(c, "请先登录")
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.Unauthorized(c, "token格式错误")
			return
		}

		claims, err := pkg.ParseToken(parts[1])
		if err != nil {
			pkg.Unauthorized(c, "token无效或已过期")
			return
		}

		// 保存用户信息到上下文
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUserName 从上下文获取用户名
func GetUserName(c *gin.Context) string {
	userName, exists := c.Get("userName")
	if !exists {
		return ""
	}
	return userName.(string)
}
