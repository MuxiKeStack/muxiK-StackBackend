package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

// 回复评论
func Reply(c *gin.Context) {
	var data model.NewCommentRequest
	if err := c.BindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("sid").(uint64)
	commentTargetId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	newCommentId, err := model.NewComment(&data, commentTargetId, false, userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	commentInfo, err := model.GetCommentInfo(newCommentId, userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, commentInfo)
}
