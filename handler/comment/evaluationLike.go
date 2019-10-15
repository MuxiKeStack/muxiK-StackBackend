package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

type evaluationLikeResponse struct {
	IsLike  bool
	LikeNum uint64
}

// 评课点赞/取消点赞
func UpdateEvaluationLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("userId").(uint64)

	err = model.UpdateEvaluationLikeState(id, userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	like := model.GetEvaluationLike(id, userId)

	data := &evaluationLikeResponse{
		IsLike:  like,
		LikeNum: 0,
	}

	handler.SendResponse(c, nil, data)
}
