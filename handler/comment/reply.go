package comment

import (
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

// 回复评论
// @Summary 回复评论
// @Tags comment
// @Param token header string true "token"
// @Param id path string true "一级评论id"
// @Param sid query string true "评论回复的目标用户的sid，若是匿名用户，则为'0'"
// @Param data body comment.newCommentRequest true "data"
// @Success 200 {object} model.CommentInfo
// @Router /comment/{id}/ [post]
func Reply(c *gin.Context) {
	var data newCommentRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)
	parentId := c.Param("id")

	// Get the user's sid whom is reply to.
	sid, ok := c.GetQuery("sid")
	if !ok {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "Target-user's sid is expected.")
		return
	}

	// Get user's id if not "0"
	var targetUserId uint32
	var err error
	if sid != "0" {
		targetUserId, err = service.GetIdBySid(sid)
		if err != nil {
			handler.SendBadRequest(c, errno.ErrGetQuery, nil, "Target-user's sid is error.")
			return
		}
	}

	// Get parentComment by id
	parentComment := &model.ParentCommentModel{Id: parentId}
	if err := parentComment.GetById(); err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, "The parent comment does not exist.")
		return
	}

	// 一级评论不匿名但没有传sid
	// if !parentComment.IsAnonymous && sid == "0" {
	// 	handler.SendBadRequest(c, errno.ErrGetQuery, nil, "Target-user's sid should not be 0!")
	// 	return
	// }

	// Check whether the targetUserId is right
	//if !(parentComment.IsAnonymous && targetUserId == 0 || !parentComment.IsAnonymous && parentComment.UserId == targetUserId) {
	//	handler.SendBadRequest(c, errno.ErrGetQuery, nil, "Sid is error, doesn't match the parentComment's")
	//	return
	//}

	// Words are limited to 300
	if len(data.Content) > 300 {
		handler.SendBadRequest(c, errno.ErrWordLimitation, nil, "Comment's content is limited to 300.")
		return
	}

	// 小程序内容安全检测
	ok, err = securityCheck.MsgSecCheck(data.Content)
	if err != nil {
		handler.SendError(c, errno.ErrSecurityCheck, nil, "check error")
		return
	} else if !ok {
		log.Errorf(err, "QQ security check msg(%s) error", data.Content)
		handler.SendBadRequest(c, errno.ErrSecurityCheck, nil, "comment content violation")
		return
	}

	var comment = &model.SubCommentModel{
		Id:           uuid.NewV4().String(),
		UserId:       userId,
		ParentId:     parentId,
		TargetUserId: targetUserId,
		Content:      data.Content,
		Time:         util.GetCurrentTime(),
		IsAnonymous:  data.IsAnonymous,
		IsValid:      true,
	}

	// Create a new subComment
	if err := comment.New(); err != nil {
		log.Error("comment.New function error. ", err)
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// Add one to the parentComment's subCommentNum
	if err := parentComment.UpdateSubCommentNum(1); err != nil {
		log.Error("UpdateSubCommentNum function error. ", err)
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// Get comment info
	commentInfo, err := service.GetSubCommentInfoById(comment.Id, userId, false)
	if err != nil {
		log.Error("service.GetSubCommentInfoById function error. ", err)
		handler.SendError(c, errno.ErrGetSubCommentInfo, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, commentInfo)

	// New message reminder
	if commentInfo.IsAnonymous {
		userId = 2 //匿名用户
	}
	err = service.NewMessageForSubComment(userId, comment, parentComment)
	if err != nil {
		log.Error("NewMessageForSubComment failed", err)
	}
}
