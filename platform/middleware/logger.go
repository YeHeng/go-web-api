package middleware

import (
	"bytes"
	"github.com/YeHeng/gtool/common/util"
	"github.com/YeHeng/gtool/platform/app"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if n, err := w.body.Write(b); err != nil {
		app.Logger.Errorf("%v", err)
		return n, err
	}
	return w.ResponseWriter.Write(b)
}

func Logger() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		buffer := util.Borrow()

		blw := &bodyLogWriter{body: buffer, ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		switch {
		case statusCode >= 400 && statusCode <= 499:
			app.Logger.Warnw(c.Errors.String(), zap.String("client_ip", clientIP),
				zap.Int("status_code", statusCode),
				zap.Duration("latency_time", latencyTime),
				zap.String("request_method", reqMethod),
				zap.String("request_uri", reqUri),
				zap.Int("response_size", c.Writer.Size()),
				zap.String("service_name", "gwebtool"),
			)
		case statusCode >= 500:
			app.Logger.Errorw(c.Errors.String(), zap.String("client_ip", clientIP),
				zap.Int("status_code", statusCode),
				zap.Duration("latency_time", latencyTime),
				zap.String("request_method", reqMethod),
				zap.String("request_uri", reqUri),
				zap.Int("response_size", c.Writer.Size()),
				zap.String("service_name", "gwebtool"),
			)
		default:

			if blw.body.Len() < 1024 {
				app.Logger.Infow(string(blw.body.Bytes()), zap.String("client_ip", clientIP),
					zap.Int("status_code", statusCode),
					zap.Duration("latency_time", latencyTime),
					zap.String("request_method", reqMethod),
					zap.String("request_uri", reqUri),
					zap.Int("response_size", c.Writer.Size()),
					zap.String("service_name", "gwebtool"),
				)
			} else {
				content := blw.body
				content.Truncate(1024)
				app.Logger.Infow(string(content.Bytes()), zap.String("client_ip", clientIP),
					zap.Int("status_code", statusCode),
					zap.Duration("latency_time", latencyTime),
					zap.String("request_method", reqMethod),
					zap.String("request_uri", reqUri),
					zap.Int("response_size", c.Writer.Size()),
					zap.String("service_name", "gwebtool"),
				)
			}
		}

		util.Return(buffer)

	}
}
