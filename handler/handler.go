package handler

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendBadRequest(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)), zap.String("cause", cause))
	c.JSON(http.StatusBadRequest, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendUnauthorized(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)), zap.String("cause", cause))
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendForbidden(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)), zap.String("cause", cause))
	c.JSON(http.StatusForbidden, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendNotFound(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)), zap.String("cause", cause))
	c.JSON(http.StatusNotFound, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendError(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Info(message, zap.String("X-Request-Id", util.GetReqID(c)), zap.String("cause", cause))
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}
