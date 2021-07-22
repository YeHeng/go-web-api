package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
)

type modifyPersonalInfoRequest struct{}

type modifyPersonalInfoResponse struct{}

// ModifyPersonalInfo 修改个人信息
// @Summary 修改个人信息
// @Description 修改个人信息
// @Tags API.admin
// @Accept json
// @Produce json
// @Param Request body modifyPersonalInfoRequest true "请求信息"
// @Success 200 {object} modifyPersonalInfoResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/modify_personal_info [patch]
func (h *handler) ModifyPersonalInfo() core.HandlerFunc {
	return func(c core.Context) {

	}
}
