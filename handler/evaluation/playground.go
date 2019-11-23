package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type playgroundResponse struct {
	Sum  int                     `json:"sum"`
	List *[]model.EvaluationInfo `json:"list"`
}

// 评课广场获取评课列表
// @Summary 评课广场获取评课列表
// @Tags comment
// @Param token header string false "游客登录则不需要此字段或为空"
// @Param limit query integer true "最大的一级评论数量"
// @Param lastId query integer true "上一次请求的最后一个评课的id，若是初始请求则为0"
// @Success 200 {object} evaluation.playgroundResponse
// @Router /evaluation/ [get]
func EvaluationPlayground(c *gin.Context) {
	log.Info("EvaluationPlayground function is called.")

	size := c.DefaultQuery("limit", "20")
	limit, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	lastIdStr := c.DefaultQuery("lastId", "0")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
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

	// 获取评课列表
	list, err := service.GetEvaluationsForPlayground(int32(lastId), int32(limit), userId, visitor)
	if err != nil {
		handler.SendError(c, errno.ErrEvaluationList, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, playgroundResponse{
		Sum:  len(*list),
		List: list,
	})
}
