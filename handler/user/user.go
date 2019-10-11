package user

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type LoginResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    model.AuthResponse `json:"data"`
}
