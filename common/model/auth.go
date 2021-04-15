package model

import (
	"gorm.io/gorm"
	"time"
)

type Credential struct {
	UserName string
	Roles    []*Role
}

type User struct {
	gorm.Model
	Email    *string `gorm:"unique;"`
	Username string  `gorm:"unique;uniqueIndex"`
	Password string
	Enabled  bool
	Expired  time.Time

	Roles []*Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name string
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
