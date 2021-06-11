package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	AddMiddleware(&swaggerMiddleware{})
}

type swaggerMiddleware struct {
}

func (m *swaggerMiddleware) Destroy() {
}

func (m *swaggerMiddleware) Init(r *gin.Engine) {
	cfg := config.Get().Feature
	if !cfg.DisableSwagger {
		fmt.Println(color.Green("* [register swagger]"))
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
	}
}
