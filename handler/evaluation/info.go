package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// 获取评课详情
// @Summary 获取评课详情
// @Tags evaluation
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param id path string true "评课id"
// @Success 200 {object} model.EvaluationInfo
// @Router /evaluation/{id}/ [get]
func GetEvaluation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	// userId获取与游客模式判断
	var userId uint32
	visitor := false

	userIdInterface, ok := c.Get("id")
	if !ok {
		visitor = true
	} else {
		userId = userIdInterface.(uint32)
		log.Info("User auth successful.")
	}

	data, err := service.GetEvaluationInfo(uint32(id), userId, visitor)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, data)
}
