package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"strings"
)

func init() {
	AddMiddleware(&corsMiddleware{})
}

type corsMiddleware struct {
}

func (m *corsMiddleware) Destroy() {
}

func (m *corsMiddleware) Init(r *gin.Engine) {

	cfg := config.Get().Feature
	if cfg.EnableCors {
		corsCfg := config.Get().Cors
		c, _ := json.Marshal(corsCfg)
		fmt.Println(color.Green(fmt.Sprintf("* [register middleware cors], options: %s", string(c))))
		r.Use(cors.New(cors.Options{
			AllowedOrigins:     strings.Split(corsCfg.AllowedOrigins, ","),
			AllowedMethods:     strings.Split(corsCfg.AllowedMethods, ","),
			AllowedHeaders:     strings.Split(corsCfg.AllowedHeaders, ","),
			AllowCredentials:   corsCfg.AllowCredentials,
			OptionsPassthrough: true,
		}))
	}
}
