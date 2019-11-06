package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type commentListResponse struct {
	ParentCommentNum  uint32                     `json:"parent_comment_num"`
	ParentCommentList *[]model.ParentCommentInfo `json:"parent_comment_list"`
}

// 获取评论列表
// @Summary 评论点赞/取消点赞
// @Tags comment
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param id path string true "评课id"
// @Param pageSize query integer true "最大的一级评论数量"
// @Param pageNum query integer true "翻页页码，默认为0"
// @Success 200 {object} comment.commentListResponse
// @Router /evaluation/{id}/comments/ [get]
func GetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
	}

	pageSize := c.DefaultQuery("pageSize", "20")
	size, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	pageNum := c.DefaultQuery("pageNum", "0")
	num, err := strconv.ParseInt(pageNum, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	// userId获取与游客模式判断
	visitor := false
	userId, ok := c.Get("id")
	if !ok {
		visitor = true
	}

	list, count, err := service.CommentList(uint32(id), int32(size), int32(num*size), userId.(uint32), visitor)
	if err != nil {
		handler.SendError(c, errno.ErrCommentList, nil, err.Error())
		return
	}

	data := commentListResponse{
		ParentCommentNum:  count,
		ParentCommentList: list,
	}

	handler.SendResponse(c, nil, data)
}
