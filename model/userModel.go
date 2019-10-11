package model

import (
	"sync"
)

type LoginModel struct {
	Sid      uint64 `json:"sid"      binding:"required"`
	Password string `json:"password" binding:"required"`
}

// User represents a registered user.
type UserModel struct {
	Sid string `json:"sid" gorm:"column:sid;not null"`
}

type UserInfo struct {
	Sid      uint64 `json:"sid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type AuthRespnse struct {
	Token string `json:"token"`
	IsNew uint8  `json:"is_new"`
}
