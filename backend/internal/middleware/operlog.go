package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"user-center/internal/model"
	"user-center/internal/repository"
)

// OperLog 操作日志中间件
func OperLog(title string, businessType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求参数
		var operParam string
		if c.Request.Method != "GET" {
			body, _ := io.ReadAll(c.Request.Body)
			operParam = string(body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 执行请求
		c.Next()

		// 异步记录日志
		go func() {
			userID := GetUserID(c)
			log := &model.SysOperLog{
				Title:         title,
				BusinessType:  businessType,
				Method:        c.HandlerName(),
				RequestMethod: c.Request.Method,
				OperUserID:    &userID,
				OperUserName:  GetUserName(c),
				OperURL:       c.Request.URL.Path,
				OperIP:        c.ClientIP(),
				OperParam:     operParam,
				Status:        "0",
				OperTime:      time.Now(),
			}
			repository.NewOperLogRepository().Create(log)
		}()
	}
}
