package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

type LoggerFormatterParams struct {
	Request      *http.Request
	TimeStamp    time.Time      // 服务器返回响应后的时间戳
	StatusCode   int            // HTTP 响应代码
	Duration     time.Duration  // 表示服务器处理请求所花费的时间
	ClientIP     string         // 客户端的 IP 地址
	Method       string         // 请求的 HTTP 方法
	Path         string         // 客户端请求的路径
	ErrorMessage string         // 在处理请求时发生的错误消息
	ResponseSize int            // 响应体的大小
	Response     string         // 响应体数据
	Query        string         // Query 请求参数
	Body         string         // Body 请求参数
	Keys         map[string]any // 在请求的上下文中设置的键值对
}

func RequestLogger() gin.HandlerFunc {
	return loggerMiddleware()
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		//bodyJson := body2string(c) // 暂不使用，会导致文件上传无法获取到文件流

		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		c.Next()

		param := LoggerFormatterParams{
			Request:      c.Request,
			Keys:         c.Keys,
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ResponseSize: c.Writer.Size(),
			Response:     writer.b.String(),
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			//Body:         bodyJson,
		}

		// 结束时间
		param.TimeStamp = time.Now()
		param.Duration = param.TimeStamp.Sub(start)

		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Path 是 /api 开头的
		if strings.HasPrefix(param.Path, "/api") {
			fmt.Println("==== API =====")
			fmt.Printf("%+v\n", param)
		}
	}
}
