package model

import (
	"gopkg.in/go-playground/validator.v9"
)

// 必须要写一个TableName函数，返回table的名字，否则gorm读取不到表。
func (u *UserModel) TableName() string {
	return "user"
}

// Update updates an user account information.
func (u *UserModel) UpdateInfo(info *UserInfoRequest) error {
	u.Avatar = info.Avatar
	u.Username = info.Username
	return DB.Self.Save(u).Error
}

// Validate the fields.
func (u *UserModel) validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// Get user info
func (u *UserModel) GetInfo() *UserInfoResponse {
	info := UserInfoResponse{
		Username: u.Username,
		Avatar:   u.Avatar,
		Sid:      u.Sid,
	}
	return &info
}

// Create creates a new user account.
func CreateUser(sid string) error {
	return DB.Self.Create(&UserModel{Sid: sid}).Error
}
