package comment

import (
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// 回复评论
// @Summary 回复评论
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "回复的评论对象id"
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

	commentTargetId := c.Param("id")

	parentId, err := model.GetParentIdByCommentTargetId(commentTargetId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	targetUserId, err := model.GetTargetUserIdByCommentTargetId(commentTargetId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var comment = &model.SubCommentModel{
		Id:           uuid.NewV4().String(),
		UserId:       userId,
		ParentId:     parentId,
		TargetUserId: targetUserId,
		Content:      data.Content,
		LikeNum:      0,
		Time:         strconv.FormatInt(time.Now().Unix(), 10),
		IsAnonymous:  data.IsAnonymous,
		IsValid:      true,
	}

	// Create a new subComment
	if err := comment.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// Get comment info
	commentInfo, err := service.GetSubCommentInfoById(comment.Id, userId, false)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentInfo)
}
