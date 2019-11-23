package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
)

type likeDataRequest struct {
	LikeState bool `json:"like_state"`
}

type likeDataResponse struct {
	LikeState bool   `json:"like_state"`
	LikeNum   uint32 `json:"like_num"`
}

// 评论点赞/取消点赞
// @Summary 评论点赞/取消点赞
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "点赞评论id"
// @Param data body comment.likeDataRequest true "当前点赞状态"
// @Success 200 {object} comment.likeDataResponse
// @Router /comment/{id}/like/ [put]
func UpdateCommentLike(c *gin.Context) {
	log.Info("UpdateCommentLike function is called.")
	var err error
	id := c.Param("id")

	// 获取请求中的点赞状态
	var bodyData likeDataRequest
	if err := c.BindJSON(&bodyData); err != nil {
		handler.SendError(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	hasLiked := model.CommentHasLiked(userId, id)

	// 判断点赞请求是否合理
	// 未点赞
	if bodyData.LikeState && !hasLiked {
		handler.SendResponse(c, errno.ErrNotLiked, nil)
		return
	}
	// 已点赞
	if !bodyData.LikeState && hasLiked {
		handler.SendResponse(c, errno.ErrHasLiked, nil)
		return
	}

	// 点赞&取消点赞
	if bodyData.LikeState {
		err = model.CommentCancelLiking(userId, id)
	} else {
		err = model.CommentLiking(userId, id)
	}

	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	data := &likeDataResponse{
		LikeState: !hasLiked,
		LikeNum:   model.GetCommentLikeSum(id),
	}

	handler.SendResponse(c, nil, data)
}
