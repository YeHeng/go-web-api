package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	AddMiddleware(&recoverMiddleware{})
}

type recoverMiddleware struct {
}

func (m *recoverMiddleware) Destroy() {
}

func (m *recoverMiddleware) Init(r *gin.Engine) {
	log := logger.Get()
	stack := config.Get().Stack
	fmt.Println(color.Green("* [register middleware recovery]"))
	r.Use(func(c *gin.Context) {
		defer func() {

			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					e := c.Error(err.(error)) // nolint: errcheck
					c.JSON(http.StatusResetContent, gin.H{
						"code":    http.StatusResetContent,
						"message": http.StatusText(http.StatusResetContent),
						"reason":  e.JSON(),
					})
					c.Abort()
					return
				}

				if stack {
					log.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": http.StatusText(http.StatusInternalServerError),
					"reason":  err,
				})
				c.Next()
			}
		}()
		c.Next()
	})
}
