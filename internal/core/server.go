package core

/*var _ Server = (*server)(nil)

type Server interface {
	// i 为了避免被其他包实现
	i()

	InitServer()

	Shutdown()
}

type server struct {
	logger *zap.Logger
	config config.Configuration
	r      *gin.Engine
}

func new(logger *zap.Logger, config config.Configuration, r *gin.Engine) Server {
	return &server{
		logger: logger,
		config: config,
		r:      r,
	}
}

func (s *server) InitServer() {

	srv = &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Logger.Infow("Server exited.")
			} else {
				logger.Logger.Fatalf("Gin start fail. %v", err)
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
	logger.Logger.Infow("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := c.WithTimeout(c.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalf("Server forced to shutdown: %v", err)
	}
}

func (s *server) Shutdown() {

	ctx, cancel := c.WithTimeout(c.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalf("Server forced to shutdown: %v", err)
	}
}

var srv *http.Server

func (s *server) i() {}*/
