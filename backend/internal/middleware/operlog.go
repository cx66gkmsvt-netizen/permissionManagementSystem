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
		// 1. 记录请求信息
		start := time.Now()
		reqMethod := c.Request.Method
		operUrl := c.Request.URL.Path
		operIp := c.ClientIP()
		method := c.HandlerName()

		var operParam string
		if reqMethod != "GET" {
			body, _ := io.ReadAll(c.Request.Body)
			operParam = string(body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 2. 获取用户信息 (JWTAuth 已执行)
		userID := GetUserID(c)
		userName := GetUserName(c)

		// 3. 执行请求
		c.Next()

		// 4. 获取响应状态
		status := "0"
		var errorMsg string
		if c.Writer.Status() != 200 {
			status = "1"
		}
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// 5. 异步记录日志 (传递捕获的变量)
		go func(uid int64, uname, m, reqM, url, ip, param, s, errMsg string, t time.Time) {
			log := &model.SysOperLog{
				Title:         title,
				BusinessType:  businessType,
				Method:        m,
				RequestMethod: reqM,
				OperUserID:    &uid,
				OperUserName:  uname,
				OperURL:       url,
				OperIP:        ip,
				OperParam:     param,
				Status:        s,
				ErrorMsg:      errMsg,
				OperTime:      t,
			}
			repository.NewOperLogRepository().Create(log)
		}(userID, userName, method, reqMethod, operUrl, operIp, operParam, status, errorMsg, start)
	}
}
