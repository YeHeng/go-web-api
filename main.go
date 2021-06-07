package main

import (
	"context"
	"errors"
	"github.com/YeHeng/go-web-api/internal/middleware"
	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/internal/pkg/plugin"
	"github.com/YeHeng/go-web-api/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Get()

	for _, p := range plugin.Get() {
		p.Init()
	}

	log := logger.Get()

	r := gin.New()

	for _, middleware := range middleware.GetMiddlewares() {
		middleware.Init(r)
	}

	r.Use()

	log.Infow("初始化Router...")
	log.Infow("开始启动APP!")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Infow("Server exited.")
			} else {
				log.Fatalf("Gin start fail. %v", err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

}
