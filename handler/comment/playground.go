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
func EvaluationPlayground(c *gin.Context) {
	pageSize := c.DefaultQuery("pageSize", "20")
	limit, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	lastIdStr := c.DefaultQuery("lastEvaluationId", "-1")
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
