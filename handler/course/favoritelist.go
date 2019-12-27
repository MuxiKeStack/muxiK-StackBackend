package course

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	//_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	_ "github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func FavoriteList(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	courseIds, err := model.GetCourseList(userId)
	if err != nil {
		log.Error("Get course list function error", err)
		return
	}

	handler.SendResponse(c, nil, courseIds)
}
