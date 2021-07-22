package middleware

import (
	"net/http"
	"time"

	"github.com/YeHeng/go-web-api/internal/code"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/YeHeng/go-web-api/pkg/errno"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	_MaxBurstSize   = 100000
	_AbortErrorName = "_abort_error_"
)

var rateLimitHandler gin.HandlerFunc

func init() {
	AddMiddleware(&rateLimitMiddleware{})
}

type rateLimitMiddleware struct {
}

func (m *rateLimitMiddleware) Init() {
	cfg := config.Get().Feature
	if cfg.EnableRate {
		limiter := rate.NewLimiter(rate.Every(time.Second*1), _MaxBurstSize)
		rateLimitHandler = func(ctx *gin.Context) {

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
		}
	}
}

func (m *rateLimitMiddleware) Apply(r *gin.Engine) {
	cfg := config.Get().Feature
	if cfg.EnableRate {
		r.Use(rateLimitHandler)
	}
}

func (m *rateLimitMiddleware) Get() gin.HandlerFunc {
	return rateLimitHandler
}

func (m *rateLimitMiddleware) Destroy() {
}
