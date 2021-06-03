package router

import (
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

// 初始化
func InitRouter(r *gin.Engine, options ...Option) {
	for _, opt := range options {
		opt(r)
	}
}
