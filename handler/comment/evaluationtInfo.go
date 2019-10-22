package comment

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"

	"github.com/gin-gonic/gin"
)

// 获取评课详情
func GetEvaluation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var userId uint32
	visitor := false
	// 游客登录&用户登录
	if t := c.Request.Header.Get("token"); len(t) == 0 {
		visitor = true
	} else {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
		}
		userId = c.MustGet("id").(uint32)
	}

	var evaluation = &model.CourseEvaluationModel{Id: uint32(id)}
	data, err := evaluation.GetInfo(userId, visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, data)
}
