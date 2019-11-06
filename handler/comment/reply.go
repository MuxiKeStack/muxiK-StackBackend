package comment

import (
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 回复评论
// @Summary 回复评论
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "评论id"
// @Param data body comment.newCommentRequest true "data"
// @Success 200 {object} model.CommentInfo
// @Router /comment/{id}/ [post]
func Reply(c *gin.Context) {
	var data newCommentRequest
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)
	commentTargetId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	parentId, err := model.GetParentIdByCommentTargetId(uint32(commentTargetId))
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	var comment = &model.CommentModel{
		UserId:          userId,
		ParentId:        parentId,
		CommentTargetId: uint32(commentTargetId),
		Content:         data.Content,
		LikeNum:         0,
		IsRoot:          false,
		Time:            strconv.FormatInt(time.Now().Unix(), 10),
		SubCommentNum:   0,
	}

	// Create a new comment
	if err := comment.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// Get comment info
	commentInfo, err := service.GetCommentInfo(comment.Id, userId, false)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentInfo)
}
