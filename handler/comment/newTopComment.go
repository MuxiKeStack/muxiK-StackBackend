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

// 新增评论请求模型
type newCommentRequest struct {
	Content     string `json:"content" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous" binding:"required"`
}

// 评论评课
// @Summary 评论评课
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "评课id"
// @Param data body comment.newCommentRequest true "data"
// @Success 200 {object} model.ParentCommentInfo
// @Router /evaluation/{id}/comment/ [post]
func CreateTopComment(c *gin.Context) {
	var data newCommentRequest
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)
	evaluationId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	var comment = &model.ParentCommentModel{
		Id:            uuid.NewV4().String(),
		UserId:        userId,
		EvaluationId:  uint32(evaluationId),
		Content:       data.Content,
		Time:          strconv.FormatInt(time.Now().Unix(), 10),
		LikeNum:       0,
		SubCommentNum: 0,
		IsAnonymous:   data.IsAnonymous,
		IsValid:       true,
	}

	// Create new comment
	if err := comment.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// Get comment info
	commentInfo, err := service.GetParentCommentInfo(comment.Id, userId, false)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentInfo)
}
