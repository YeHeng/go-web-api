package main

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
	"github.com/YeHeng/go-web-api/internal/router"
	middleware2 "github.com/YeHeng/go-web-api/internal/router/middleware"
	"github.com/YeHeng/go-web-api/pkg"
	config2 "github.com/YeHeng/go-web-api/pkg/config"
	"github.com/YeHeng/go-web-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	config2.LoadConfig()
	logger.InitLogger()
	middleware2.InitDb()
	pkg.InitCasbin()

	r := gin.New()
	r.Use(middleware2.Logger(), middleware2.Recovery(false))
	middleware2.InitJwt(r)
	router.InitRouter(r)
	logger.Logger.Infow("初始化Router...")
	logger.Logger.Infow("开始启动APP!")

	config := config2.Config

	core.InitServer(config, r)

}
