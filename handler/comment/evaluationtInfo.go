package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 获取评课详情
func GetEvaluation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	// userId获取与游客模式判断
	visitor := false
	userId, ok := c.Get("id")
	if !ok {
		visitor = true
	}

	data, err := service.GetEvaluationInfo(uint32(id), userId.(uint32), visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, data)
}
