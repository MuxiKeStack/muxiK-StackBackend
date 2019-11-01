package comment

import (
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

// 回复评论
func Reply(c *gin.Context) {
	var data model.NewCommentRequest
	if err := c.BindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("id").(uint32)
	commentTargetId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	parentId, err := model.GetParentIdByCommentTargetId(uint32(commentTargetId))
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
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

	if err := comment.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	commentInfo, err := comment.GetInfo(userId, false)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, commentInfo)
}
