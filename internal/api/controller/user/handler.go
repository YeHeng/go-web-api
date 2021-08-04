package user

import (
	"github.com/YeHeng/go-web-api/internal/pkg/context"
)

var _ Handler = (*handler)(nil)

type handler struct {
}

type Handler interface {
	i()

	// ModifyPassword 修改密码
	// @Tags API.admin
	// @Router /api/admin/modify_password [patch]
	ModifyPassword() context.HandlerFunc

	// Detail 个人信息
	// @Tags API.admin
	// @Router /api/admin/info [get]
	Detail() context.HandlerFunc

	// ModifyPersonalInfo 修改个人信息
	// @Tags API.admin
	// @Router /api/admin/modify_personal_info [patch]
	ModifyPersonalInfo() context.HandlerFunc

	// Create 新增管理员
	// @Tags API.admin
	// @Router /api/admin [post]
	Create() context.HandlerFunc

	// List 管理员列表
	// @Tags API.admin
	// @Router /api/admin [get]
	List() context.HandlerFunc

	// Delete 删除管理员
	// @Tags API.admin
	// @Router /api/admin/{id} [delete]
	Delete() context.HandlerFunc

	// ResetPassword 重置密码
	// @Tags API.admin
	// @Router /api/admin/reset_password/{id} [patch]
	ResetPassword() context.HandlerFunc
}

func (h *handler) i() {}
