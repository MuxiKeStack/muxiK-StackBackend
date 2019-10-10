package model

import (
	"sync"
)

type BaseModel struct {
	Sid        uint64     `gorm:"primary_key;column:sid" json:"sid"`
}

type UserInfo struct {
	Sid	       	uint64 `json:"sid"`
	Username  	string `json:"username"`
	Avatar		string `json:"avatar"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}
