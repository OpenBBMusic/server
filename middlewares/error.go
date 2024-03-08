package middlewares

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggerFormatterParams struct {
	StatusCode   int            // HTTP 响应代码
	Duration     time.Duration  // 表示服务器处理请求所花费的时间
	ClientIP     string         // 客户端的 IP 地址
	Method       string         // 请求的 HTTP 方法
	ErrorMessage string         // 在处理请求时发生的错误消息
	ResponseSize int            // 响应体的大小
	Query        string         // Query 请求参数
	Keys         map[string]any // 在请求的上下文中设置的键值对
}

func RequestLogger() gin.HandlerFunc {
	return loggerMiddleware()
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		// 开始时间
		start := time.Now()
		if strings.HasPrefix(reqPath, "/api") {
			log.Println("API START", reqPath)
		}

		c.Next()

		param := LoggerFormatterParams{
			Keys:         c.Keys,
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ResponseSize: c.Writer.Size(),
			Query:        c.Request.URL.RawQuery,
		}

		// 结束时间
		param.Duration = time.Since(start)
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Path 是 /api 开头的
		if strings.HasPrefix(reqPath, "/api") {
			log.Printf("API END %v %+v\n", reqPath, param)
		}
	}
}
