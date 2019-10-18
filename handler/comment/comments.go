package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

type commentListResponse struct {
	ParentCommentNum  uint64                     `json:"parent_comment_num"`
	ParentCommentList *[]model.ParentCommentInfo `json:"parent_comment_list"`
}

// 获取评论列表
func GetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	pageSize := c.DefaultQuery("pageSize", "20")
	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	} else if size <= 0 {
		handler.SendBadRequest(c, err, nil, "PageSize error")
	}

	lastIdStr := c.DefaultQuery("lastCommentId", "-1")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 64)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var userId uint64
	visitor := false
	// 游客登录
	if t := c.Request.Header.Get("token"); len(t) == 0 {
		visitor = true
	} else {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
		}
		userId = c.MustGet("sid").(uint64)
	}

	list, count, err := model.GetCommentList(id, lastId, size, userId, visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	data := commentListResponse{
		ParentCommentNum:  count,
		ParentCommentList: list,
	}

	handler.SendResponse(c, nil, data)
}
