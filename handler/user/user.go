package user

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type CreateRequest struct {
	model.LoginModel
}

type LoginResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    model.AuthResponse `json:"data"`
}
