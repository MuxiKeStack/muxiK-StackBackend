package model

import (
	"gopkg.in/go-playground/validator.v9"
)

func (u *UserModel) TableName() string {
	return "user"
}

// Update updates an user account information.
func (u *UserModel) Update() error {
	return DB.Self.Save(u).Error
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
func GetUserById(id uint64) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("id = ?", id).First(&u)
	return u, d.Error
}

// Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
