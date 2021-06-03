package main

import (
	"github.com/YeHeng/gtool/pkg"
	"github.com/YeHeng/gtool/platform/app"
	"github.com/YeHeng/gtool/platform/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	app.LoadConfig()
	app.InitLogger()
	middleware.InitDb()
	pkg.InitCasbin()

	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery(false))
	middleware.InitJwt(r)
	app.InitRouter(r)
	app.Logger.Infow("初始化Router...")
	app.Logger.Infow("开始启动APP!")

	config := app.Config

	app.InitServer(config, r)

}
