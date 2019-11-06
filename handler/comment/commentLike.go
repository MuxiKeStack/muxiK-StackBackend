package comment

import (
	"errors"
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

	var comment = &model.CommentModel{Id: uint32(id)}
	hasLiked := comment.HasLiked(userId)

	// 取消点赞
	if bodyData.IsLike {
		if !hasLiked {
			err = errors.New("Has not liked yet. ")
			handler.SendResponse(c, err, nil)
			return
		}
		err = comment.CancelLiking(userId)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
		err = comment.UpdateLikeNum(-1)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
	} else {
		// 点赞

		if hasLiked {
			err = errors.New("Has already liked. ")
			handler.SendResponse(c, err, nil)
			return
		}
		err = comment.Like(userId)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
		err = comment.UpdateLikeNum(1)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
	}

	data := &likeDataResponse{
		IsLike:  !hasLiked,
		LikeNum: comment.LikeNum,
	}

	handler.SendResponse(c, nil, data)
}
