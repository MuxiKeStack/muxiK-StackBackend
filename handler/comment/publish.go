package comment

import (
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 发布评课的请求数据
type evaluationPublishRequest struct {
	CourseId            string  `json:"course_id" binding:"required"`
	CourseName          string  `json:"course_name" binding:"required"`
	Rate                float32 `json:"rate" binding:"required"`
	AttendanceCheckType uint8   `json:"attendance_check_type" binding:"required"`
	ExamCheckType       uint8   `json:"exam_check_type" binding:"required"`
	Content             string  `json:"content" binding:"required"`
	IsAnonymous         bool    `json:"is_anonymous" binding:"required"`
	Tags                []uint8 `json:"tags" binding:"required"`
}

type evaluationPublishResponse struct {
	EvaluationId uint32 `json:"evaluation_id"`
}

// Publish a new  course evaluation.
// @Summary 发布评课
// @Tags comment
// @Param token header string true "token"
// @Param data body comment.evaluationPublishRequest true "data"
// @Success 200 {object} comment.evaluationPublishResponse
// @Router /evaluation/ [post]
func Publish(c *gin.Context) {
	var data evaluationPublishRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	var evaluation = &model.CourseEvaluationModel{
		CourseId:            data.CourseId,
		CourseName:          data.CourseName,
		UserId:              userId,
		Rate:                data.Rate,
		AttendanceCheckType: data.AttendanceCheckType,
		ExamCheckType:       data.ExamCheckType,
		Content:             data.Content,
		LikeNum:             0,
		CommentNum:          0,
		Tags:                service.TagArrayToStr(data.Tags),
		IsAnonymous:         data.IsAnonymous,
		IsValid:             true,
		Time:                strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err := evaluation.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 更新数据库中课程的评分信息
	if err := model.UpdateCourseRateByEvaluation(evaluation.CourseId, data.Rate); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &evaluationPublishResponse{EvaluationId: evaluation.Id})
}
