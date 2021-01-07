package comment

import (
	"strconv"
	"unicode/utf8"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/MuxiKeStack/muxiK-StackBackend/util/securityCheck"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	uuid "github.com/satori/go.uuid"
)

// 新增评论请求模型
type newCommentRequest struct {
	Content     string `json:"content" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous" binding:"-"`
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
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)
	evaluationId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	// Words are limited to 300
	if utf8.RuneCountInString(data.Content) > 300 {
		handler.SendBadRequest(c, errno.ErrWordLimitation, nil, "Comment's content is limited to 300.")
		return
	}

	// 小程序内容安全检测
	// 问题：历史token已过兼容期
	// 先关掉，之后查验
	ok, err := securityCheck.MsgSecCheck(data.Content)
	if err != nil {
		log.Error("QQ security check function error", err)
		// handler.SendError(c, errno.ErrSecurityCheck, nil, "check error")
		// return
	} else if !ok {
		log.Errorf(err, "QQ security check msg(%s) error", data.Content)
		// handler.SendBadRequest(c, errno.ErrSecurityCheck, nil, "comment content violation")
		// return
	}

	var comment = &model.ParentCommentModel{
		Id:            uuid.NewV4().String(),
		UserId:        userId,
		EvaluationId:  uint32(evaluationId),
		Content:       data.Content,
		Time:          util.GetCurrentTime(),
		SubCommentNum: 0,
		IsAnonymous:   data.IsAnonymous,
		IsValid:       true,
	}

	// Create new comment
	if err := comment.New(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// Add one to the evaluation's comment sum
	evaluation := &model.CourseEvaluationModel{Id: uint32(evaluationId)}
	if err := evaluation.GetById(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	if err := evaluation.UpdateCommentNum(1); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// Get comment info
	commentInfo, err := service.GetParentCommentInfo(comment.Id, userId, false)
	if err != nil {
		handler.SendError(c, errno.ErrGetParentCommentInfo, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentInfo)

	// New message reminder
	if commentInfo.IsAnonymous {
		userId = 2 //匿名用户
	}

	err = service.NewMessageForParentComment(userId, comment, evaluation)
	if err != nil {
		log.Error("NewMessageForParentComment failed", err)
	}
}
