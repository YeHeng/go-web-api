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

var swaggerHandler gin.HandlerFunc

type swaggerMiddleware struct {
}

func (m *swaggerMiddleware) Get() gin.HandlerFunc {
	return swaggerHandler
}

func (m *swaggerMiddleware) Init() {
}

func (m *swaggerMiddleware) Apply(r *gin.Engine) {
	cfg := config.Get().Feature
	if !cfg.DisableSwagger {
		swaggerHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
		fmt.Println(color.Green("* [register swagger]"))
		r.GET("/swagger/*any", swaggerHandler) // register swagger
	}
}

func (m *swaggerMiddleware) Destroy() {
}
