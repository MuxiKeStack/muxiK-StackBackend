package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

type likeDataResponse struct {
	IsLike  bool   `json:"is_like"`
	LikeNum uint32 `json:"like_num"`
}

type likeDataRequest struct {
	IsLike bool `json:"is_like"`
}

// 评课点赞/取消点赞
// @Summary 评论点赞/取消点赞
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "点赞的评课id"
// @Param data body comment.likeDataRequest true "data"
// @Success 200 {object} comment.likeDataResponse
// @Router /evaluation/{id}/like/ [put]
func UpdateEvaluationLike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 获取请求中当前的点赞状态
	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	var evaluation = &model.CourseEvaluationModel{Id: uint32(id)}
	hasLiked := evaluation.HasLiked(userId)

	// 判断点赞请求是否合理
	// 未点赞
	if bodyData.IsLike && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}
	//	已点赞
	if !bodyData.IsLike && hasLiked {
		handler.SendResponse(c, errno.ErrHasLiked, nil)
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
	num := 1
	if bodyData.IsLike {
		num = -1
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
