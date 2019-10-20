package comment

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

type responseData struct {
	EvaluationId uint32 `json:"evaluation_id"`
}

// 发布评课
func Publish(c *gin.Context) {
	var data model.EvaluationPublish
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	userId := c.MustGet("id").(uint32)
	evaluationId, err := model.NewEvaluation(&data, userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	// 更新数据库中课程的评分信息
	if err := model.UpdateCourseByEva(evaluationId, data.Rate); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, &responseData{EvaluationId: evaluationId})
}
