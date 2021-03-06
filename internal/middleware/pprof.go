package middleware

import (
	"fmt"

	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func init() {
	AddMiddleware(&pprofMiddleware{})
}

type pprofMiddleware struct {
}

func (m *pprofMiddleware) Init() {
}

func (m *pprofMiddleware) Apply(r *gin.Engine) {
	cfg := config.Get().Feature
	if !cfg.DisablePProf {
		fmt.Println(color.Green("* [register middleware pprof]"))
		pprof.Register(r) // register pprof to gin
	}
}

func (m *pprofMiddleware) Get() gin.HandlerFunc {
	return nil
}

func (m *pprofMiddleware) Destroy() {
}
