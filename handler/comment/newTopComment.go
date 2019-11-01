package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 新增评论请求模型
type NewCommentRequest struct {
	Content     string `json:"content"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// 评论评课
func CreateTopComment(c *gin.Context) {
	var data NewCommentRequest
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

	var comment = &model.CommentModel{
		UserId:          userId,
		ParentId:        0,
		CommentTargetId: uint32(evaluationId),
		Content:         data.Content,
		LikeNum:         0,
		IsRoot:          true,
		Time:            strconv.FormatInt(time.Now().Unix(), 10),
		SubCommentNum:   0,
	}

	// Create new comment
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
