package service

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/model"
)

// HaveUser determines whether there is this user or not by the user identifier.
func HaveUser(sid string) uint8 {
	var num int
	DB.Self.Model(&UserModel{}).Where("sid = ?", sid).Count(&num)
	if num == 0 {
		return 1
	}
	return 0
}

// GetIdBySid gets a user id by user's sid
func GetIdBySid(sid string) (uint32, error) {
	u := &UserModel{}
	d := DB.Self.Where("sid = ?", sid).First(&u)
	return u.Id, d.Error
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
func GetUserInfoById(id uint32) (*UserInfoResponse, error) {
	u, err := GetUserById(id)
	if err != nil {
		return &UserInfoResponse{}, err
	}
	info := u.GetInfo()
	return info, nil
}

func GetUserInfoRById(id uint32) (*UserInfoRequest, error) {
	u, err := GetUserById(id)
	if err != nil {
		return &UserInfoRequest{}, err
	}
	info := u.GetInfo()

	return &UserInfoRequest{
		Username: info.Username,
		Avatar:   info.Avatar,
	}, nil
}

// UpdateInfoById update user information by Id
func UpdateInfoById(id uint32, info *UserInfoRequest) error {
	u, err := GetUserById(id)
	if err != nil {
		return err
	}
	if err = u.UpdateInfo(info); err != nil {
		return err
	}
	return nil
}
