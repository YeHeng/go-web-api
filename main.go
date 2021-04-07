package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/YeHeng/gtool/app"
	"github.com/YeHeng/gtool/app/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery(false) /*, gzip.Gzip(gzip.DefaultCompression)*/)
	app.Logger.Infow("初始化Router...")
	r.Static("/static", "./web")
	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))
	r.NoRoute(noRoute)
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(".", "client", "dist", "index.html"))
	})

	app.Logger.Infof("开始启动APP!")

	config := app.Config

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				app.Logger.Infow("Server exited.")
			} else {
				app.Logger.Fatalf("Gin start fail. %v", err)
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
	app.Logger.Infow("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

}

func noRoute(c *gin.Context) {
	path := strings.Split(c.Request.URL.Path, "/")
	if path[1] == "api" {
		c.JSON(http.StatusNotFound, gin.H{"msg": "no route", "body": nil})
	} else {
		c.HTML(http.StatusOK, "index.html", "")
	}
}
