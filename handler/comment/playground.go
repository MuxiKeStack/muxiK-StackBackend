package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type playgroundResponse struct {
	Sum  int                     `json:"sum"`
	List *[]model.EvaluationInfo `json:"list"`
}

// 评课广场获取评课列表
// @Summary 评课广场获取评课列表
// @Tags comment
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param pageSize query integer true "最大的一级评论数量"
// @Param lastEvaluationId query integer true "上一次请求的最后一个评课的id，若是初始请求则为空或-1"
// @Success 200 {object} comment.playgroundResponse
// @Router /evaluation/list/ [get]
func EvaluationPlayground(c *gin.Context) {
	pageSize := c.DefaultQuery("pageSize", "20")
	limit, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	lastIdStr := c.DefaultQuery("lastEvaluationId", "")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
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

	// 获取评课列表
	list, err := service.EvaluationList(int32(lastId), int32(limit), userId.(uint32), visitor)
	if err != nil {
		handler.SendError(c, errno.ErrEvaluationList, nil, err.Error())
		return
	}

	data := playgroundResponse{
		Sum:  len(*list),
		List: list,
	}

	handler.SendResponse(c, nil, data)
}
