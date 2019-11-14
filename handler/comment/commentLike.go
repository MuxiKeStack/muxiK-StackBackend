package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 评论点赞/取消点赞
// @Summary 评论点赞/取消点赞
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "点赞评论id"
// @Param data body comment.likeDataRequest true "data"
// @Success 200 {object} comment.likeDataResponse
// @Router /comment/{id}/like/ [put]
func UpdateCommentLike(c *gin.Context) {
	var err error

	id := c.Param("id")

	// 获取请求中当前的点赞状态
	var bodyData likeDataRequest
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	hasLiked := model.CommentHasLiked(userId, id)

	// 判断点赞请求是否合理
	// 未点赞
	if bodyData.IsLike && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}
	// 已点赞
	if !bodyData.IsLike && hasLiked {
		handler.SendResponse(c, errno.ErrHasLiked, nil)
		return
	}

	// 点赞&取消点赞
	if bodyData.IsLike {
		err = model.CommentCancelLiking(userId, id);
	} else {
		err = model.CommentLike(userId, id)
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
	count, err := service.UpdateCommentLikeNum(id, num)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	data := &likeDataResponse{
		IsLike:  !hasLiked,
		LikeNum: count,
	}

	handler.SendResponse(c, nil, data)
}
