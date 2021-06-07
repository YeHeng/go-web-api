package middleware

import (
	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func init() {
	AddMiddleware(&pprofMiddleware{})
}

type pprofMiddleware struct {
}

func (m *pprofMiddleware) Destroy() {
}

func (m *pprofMiddleware) Init(r *gin.Engine) {
	cfg := config.Get().Feature
	if !cfg.DisablePProf {
		pprof.Register(r) // register pprof to gin
		logger.Get().Infow("* [register pprof]")
	}
}
