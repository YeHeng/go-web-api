package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Init()
	Apply(r *gin.Engine)
	Get() gin.HandlerFunc
	Destroy()
}

var plugins = make([]Middleware, 0)

func AddMiddleware(middleware Middleware) {
	plugins = append(plugins, middleware)
}

func GetMiddlewares() []Middleware {
	return plugins
}
