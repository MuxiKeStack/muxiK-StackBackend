package evaluation

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type historyEvaluationsResponse struct {
	Sum      int                     `json:"sum"`
	List     *[]model.EvaluationInfo `json:"list"`
}

// @Summary 个人历史评课
// @Tags evaluation
// @Param token header string true "token"
// @Param limit query integer true "评课数"
// @Param last_id query integer true "上一次请求的最后一个评课的id，若是初始请求则为0"
// @Success 200 {object} historyEvaluationsResponse
// @Router /user/evaluations/ [get]
func GetHistoryEvaluations(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "limit parse error.")
		return
	}

	lastIdStr := c.DefaultQuery("last_id", "0")
	lastId, err := strconv.ParseInt(lastIdStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}

	list, err := service.GetHistoryEvaluationsByUserId(userId, int32(lastId), int32(limit))
	if err != nil {
		handler.SendError(c, errno.ErrGetHistoryEvaluations, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, historyEvaluationsResponse{
		Sum:      len(*list),
		List:     list,
	})
}
