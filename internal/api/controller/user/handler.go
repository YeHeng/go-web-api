package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
)

var _ Handler = (*handler)(nil)

type handler struct {
}

type Handler interface {
	i()

	// ModifyPassword 修改密码
	// @Tags API.admin
	// @Router /api/admin/modify_password [patch]
	ModifyPassword() core.HandlerFunc

	// Detail 个人信息
	// @Tags API.admin
	// @Router /api/admin/info [get]
	Detail() core.HandlerFunc

	// ModifyPersonalInfo 修改个人信息
	// @Tags API.admin
	// @Router /api/admin/modify_personal_info [patch]
	ModifyPersonalInfo() core.HandlerFunc

	// Create 新增管理员
	// @Tags API.admin
	// @Router /api/admin [post]
	Create() core.HandlerFunc

	// List 管理员列表
	// @Tags API.admin
	// @Router /api/admin [get]
	List() core.HandlerFunc

	// Delete 删除管理员
	// @Tags API.admin
	// @Router /api/admin/{id} [delete]
	Delete() core.HandlerFunc

	// ResetPassword 重置密码
	// @Tags API.admin
	// @Router /api/admin/reset_password/{id} [patch]
	ResetPassword() core.HandlerFunc
}

func (h *handler) i() {}
