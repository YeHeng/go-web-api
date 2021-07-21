package user_repo

import "time"

// User 用户表
//go:generate gormgen -structs User -input .
type User struct {
	Id          int64     //
	Username    string    // 用户名
	Password    string    //
	Mobile      string    //
	IsDeleted   int32     //
	CreatedTime time.Time `gorm:"time"` //
	UpdateTime  time.Time `gorm:"time"` //
}
