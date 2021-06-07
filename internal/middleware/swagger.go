package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/env"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	Add(&swaggerMiddleware{})
}

type swaggerMiddleware struct {
}

func (m *swaggerMiddleware) Apply(r *gin.Engine) {
	if !env.Active().IsPro() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
		fmt.Println(color.Green("* [register swagger]"))
	}
}
