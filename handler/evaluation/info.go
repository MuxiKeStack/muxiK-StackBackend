package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/MuxiKeStack/muxiK-StackBackend/log"
	"github.com/gin-gonic/gin"
)

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

	// 从数据库中获取评课记录
	evaluation := model.NewEvaluation(uint32(id))
	if err := evaluation.GetById(); err != nil {
		log.Error("evaluation.GetById function error.", err)
		handler.SendError(c, errno.ErrGetEvaluationInfo, nil, err.Error())
		return
	}

	data, err := service.GetEvaluationInfo(evaluation, userId, visitor)
	if err != nil {
		handler.SendError(c, errno.ErrGetEvaluationInfo, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, data)
}
