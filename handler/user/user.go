package user

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type CreateLoginRequest struct {
	model.LoginModel
}

type CreatePostInfoRequest struct {
	model.UserInfo
}

type LoginResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    model.AuthResponse `json:"data"`
}

type InfoResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    model.UserInfo `json:"data"`
}
