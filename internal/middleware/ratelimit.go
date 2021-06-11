package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/internal/code"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/YeHeng/go-web-api/pkg/errno"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	_MaxBurstSize   = 100000
	_AbortErrorName = "_abort_error_"
)

func init() {
	AddMiddleware(&rateLimitMiddleware{})
}

type rateLimitMiddleware struct {
}

func (m *rateLimitMiddleware) Destroy() {
}

func (m *rateLimitMiddleware) Init(r *gin.Engine) {
	cfg := config.Get().Feature
	if cfg.EnableRate {
		fmt.Println(color.Green("* [register swagger]"))
		limiter := rate.NewLimiter(rate.Every(time.Second*1), _MaxBurstSize)
		r.Use(func(ctx *gin.Context) {

			if !limiter.Allow() {

				err := errno.NewError(
					http.StatusTooManyRequests,
					code.TooManyRequests,
					code.Text(code.TooManyRequests))

				httpCode := err.GetHttpCode()
				if httpCode == 0 {
					httpCode = http.StatusInternalServerError
				}

				ctx.AbortWithStatus(httpCode)
				ctx.Set(_AbortErrorName, err)

				return
			}

			ctx.Next()
		})
	}
}
