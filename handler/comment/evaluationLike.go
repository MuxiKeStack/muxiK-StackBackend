package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

type likeDataResponse struct {
	IsLike  bool
	LikeNum uint32
}

type likeDataRequest struct {
	IsLike bool
}

// 评课点赞/取消点赞
func UpdateEvaluationLike(c *gin.Context) {
	// 获取请求中当前的点赞状态
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	var evaluation = &model.CourseEvaluationModel{Id: uint32(id)}
	hasLiked := evaluation.HasLiked(userId)

	if bodyData.IsLike && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}

	// 点赞或者取消点赞
	if bodyData.IsLike {
		err = evaluation.CancelLiking(userId)
	} else {
		err = evaluation.Like(userId)
	}

	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 更新点赞数
	num := -1
	if bodyData.IsLike {
		num = 1
	}
	err = evaluation.UpdateLikeNum(num)

	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	data := &likeDataResponse{
		IsLike:  !hasLiked,
		LikeNum: evaluation.LikeNum,
	}

	handler.SendResponse(c, nil, data)
}
