package model

import (
	"gopkg.in/go-playground/validator.v9"
)

func (u *UserModel) TableName() string {
	return "user"
}

// Update updates an user account information.
func (u *UserModel) updateInfo(info *UserInfo) error {
	u.Avatar = info.Avatar
	u.Username = info.Username
	return DB.Self.Save(u).Error
}

// Get user info
func (u *UserModel) getInfo() *UserInfo {
	info := UserInfo{
		Username: u.Username,
		Avatar:   u.Avatar,
	}
	return &info
}

// Validate the fields.
func (u *UserModel) validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// Create creates a new user account.
func CreateUser(sid string) error {
	return DB.Self.Create(&UserModel{Sid: sid}).Error
}

// HaveUser determines whether there is this user or not by the user identifier.
func HaveUser(sid string) (uint8, error) {
	var num int
	DB.Self.Model(&UserModel{}).Where("sid = ?", sid).Count(num)
	if num == 0 {
		return 0, nil
	}
	return 1, nil
}

// GetUser gets an user by the student identifier.
func GetUserBySid(sid string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("sid = ?", sid).First(&u)
	return u, d.Error
}

// GetUser gets an user by the user identifier.
func GetUserById(id uint32) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("id = ?", id).First(&u)
	return u, d.Error
}

// GetUserInfoById gets user information by userId.
func GetUserInfoById(id uint32) (*UserInfo, error) {
	u, err := GetUserById(id)
	if err != nil {
		return &UserInfo{}, err
	}
	info := u.getInfo()
	return info, nil
}

// GetUserInfoBySid gets user information by Sid
func GetUserInfoBySid(sid string) (*UserInfo, error) {
	u, err := GetUserBySid(sid)
	if err != nil {
		return &UserInfo{}, nil
	}
	info := u.getInfo()
	return info, nil
}

// UpdateInfoById update user information by Id
func UpdateInfoById(id uint32, info *UserInfo) error {
	u, err := GetUserById(id)
	if err != nil {
		return err
	}
	if err = u.updateInfo(info); err != nil {
		return err
	}
	return nil
}
