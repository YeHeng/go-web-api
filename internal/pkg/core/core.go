package core

import (
	c "context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YeHeng/go-web-api/internal/middleware"
	"github.com/YeHeng/go-web-api/internal/pkg/factory"
	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var _ Mux = (*mux)(nil)

// Mux http mux
type Mux interface {
	http.Handler
}

type mux struct {
	engine *gin.Engine
}

func (m mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.engine.ServeHTTP(writer, request)
}

func Create() (Mux, error) {

	log := logger.Get()
	if log == nil {
		return nil, errors.New("logger required")
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()

	mux := &mux{
		engine: gin.New(),
	}

	dispatcher(mux.engine)

	for _, m := range middleware.GetMiddlewares() {
		m.Init()
		m.Apply(mux.engine)
	}

	mux.engine.NoMethod(wrapHandlers(DisableTrace, func(c Context) {
		c.GetContext().JSON(http.StatusMethodNotAllowed, gin.H{
			"code":    http.StatusMethodNotAllowed,
			"message": http.StatusText(http.StatusMethodNotAllowed),
			"uri":     c.URI(),
		})
	})...)
	mux.engine.NoRoute(wrapHandlers(DisableTrace, func(c Context) {
		c.GetContext().JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": http.StatusText(http.StatusNotFound),
			"uri":     c.URI(),
		})
	})...)

	cfg := config.Get()
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux.engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {

			for _, m := range middleware.GetMiddlewares() {
				m.Destroy()
			}

			for _, p := range factory.GetAllBeans() {
				p.Destroy()
			}

			if errors.Is(err, http.ErrServerClosed) {
				log.Info("Server exited.")
			} else {
				log.Fatal(fmt.Sprintf("Gin start fail. %v", err))
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
	log.Info("Shutting down server...")

	for _, m := range middleware.GetMiddlewares() {
		m.Destroy()
	}

	for _, p := range factory.GetAllBeans() {
		p.Destroy()
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := c.WithTimeout(c.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	return mux, nil
}

func DisableTrace(ctx Context) {
	ctx.disableTrace()
}

func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {

			ctx := newContext(c)
			defer releaseContext(ctx)

			handler(ctx)
		}
	}

	return funcs
}
