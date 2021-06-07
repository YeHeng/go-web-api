package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func init() {
	Add(&pprofMiddleware{})
}

type pprofMiddleware struct {
}

func (m *pprofMiddleware) Apply(r *gin.Engine) {
	cfg := config.Get().Feature
	if !cfg.DisablePProf {
		pprof.Register(r) // register pprof to gin
		fmt.Println(color.Green("* [register pprof]"))
	}
}
