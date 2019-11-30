package middleware

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		ctx, err := token.ParseRequest(c)
		if err != nil {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, err.Error())
			c.Abort()
			return
		}
		//gin.Context内有一个keys是一个map[string][interface{}]
		c.Set("id", ctx.Id)
		c.Next()
	}
}

func VisitorAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		ctx, err := token.ParseRequest(c)
		if err != nil {
			log.Info("Token is invalid. Entry visitor mode.")
		} else {
			c.Set("id", ctx.Id)
		}

		c.Next()
	}
}
