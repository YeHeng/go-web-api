package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
)

type deleteRequest struct{}

type deleteResponse struct{}

// Delete 删除管理员
// @Summary 删除管理员
// @Description 删除管理员
// @Tags API.admin
// @Accept json
// @Produce json
// @Param Request body deleteRequest true "请求信息"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/{id} [delete]
func (h *handler) Delete() core.HandlerFunc {
	return func(c core.Context) {

	}
}
