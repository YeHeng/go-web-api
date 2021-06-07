package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Apply(r *gin.Engine)
}

var middlewares = make([]Middleware, 0)

func Add(middleware Middleware) {
	middlewares = append(middlewares, middleware)
}

func Get() []Middleware {
	return middlewares
}
