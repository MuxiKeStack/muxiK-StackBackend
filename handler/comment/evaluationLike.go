package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

type commentLikeResponse struct {
	IsLike  bool
	LikeNum uint64
}

type likeResponse struct {
	IsLike bool
}

// 评课点赞/取消点赞
func UpdateEvaluationLike(c *gin.Context) {
	// 获取请求中当前的点赞状态
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var d likeResponse
	if err := c.BindJSON(&d); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("sid").(uint64)

	err = model.UpdateEvaluationLikeState(id, userId, d.IsLike)

	data := &commentLikeResponse{
		IsLike:  model.GetEvaluationLikeState(id, userId),
		LikeNum: model.GetEvaluationLikeNum(id),
	}

	// 数据库的点赞状态和请求的点赞状态相冲突
	if err != nil {
		handler.SendResponse(c, err, data)
	}

	handler.SendResponse(c, nil, data)
}
