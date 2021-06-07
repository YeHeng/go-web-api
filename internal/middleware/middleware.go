package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Init(r *gin.Engine)
	Destroy()
}

var plugins = make([]Middleware, 0)

func AddMiddleware(middleware Middleware) {
	plugins = append(plugins, middleware)
}

func GetMiddlewares() []Middleware {
	return plugins
}
