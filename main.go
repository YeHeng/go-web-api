package main

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
	"github.com/YeHeng/go-web-api/internal/router"
	"github.com/YeHeng/go-web-api/internal/router/middleware"
	"github.com/YeHeng/go-web-api/pkg"
	"github.com/YeHeng/go-web-api/pkg/config"
	db2 "github.com/YeHeng/go-web-api/pkg/db"
	"github.com/YeHeng/go-web-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()
	logger.InitLogger()
	db2.InitDb()
	pkg.InitCasbin()

	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery(false))
	middleware.InitJwt(r)
	router.InitRouter(r)
	logger.Logger.Infow("初始化Router...")
	logger.Logger.Infow("开始启动APP!")

	config := config.Config

	core.InitServer(config, r)

}
