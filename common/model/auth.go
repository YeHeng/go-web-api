package model

import (
	"time"

	"gorm.io/gorm"
)

type Credential struct {
	UserName string
	Roles    []*Role
}

type User struct {
	ID          int64     `gorm:"primaryKey;"`
	Email       string    `gorm:"uniqueIndex;size:512;not null;"`
	Username    string    `gorm:"uniqueIndex;size:512;not null;"`
	Password    string    `gorm:"not null"`
	Enabled     bool      `gorm:"not null;default:false"`
	Expired     time.Time `gorm:"not null;"`
	CreatedTime time.Time `gorm:"autoCreateTime;not null;"`
	UpdateTime  time.Time `gorm:"autoUpdateTime;not null;"`
}

type Role struct {
	gorm.Model
	Name string
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
