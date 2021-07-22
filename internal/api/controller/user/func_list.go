package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
)

type listRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type listResponse struct{}

// List 管理员列表
// @Summary 管理员列表
// @Description 管理员列表
// @Tags API.admin
// @Accept json
// @Produce json
// @Param Request body listRequest true "请求信息"
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin [get]
func (h *handler) List() core.HandlerFunc {
	return func(c core.Context) {

	}
}
