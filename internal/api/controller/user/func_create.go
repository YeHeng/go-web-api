package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
)

type createRequest struct{}

type createResponse struct{}

// Create 新增管理员
// @Summary 新增管理员
// @Description 新增管理员
// @Tags API.admin
// @Accept json
// @Produce json
// @Param Request body createRequest true "请求信息"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin [post]
func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {

	}
}
