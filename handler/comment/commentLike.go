package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

// 评论点赞/取消点赞
func UpdateCommentLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var d likeResponse
	if err := c.BindJSON(&d); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("id").(uint32)

	err = model.UpdateCommentLikeState(uint32(id), userId, d.IsLike)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	data := &commentLikeResponse{
		IsLike:  model.GetCommentLikeState(uint32(id), userId),
		LikeNum: model.GetCommentLikeNum(uint32(id)),
	}

	handler.SendResponse(c, nil, data)
}
