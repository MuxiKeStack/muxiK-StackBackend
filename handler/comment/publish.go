package comment

import (
	"strconv"
	"time"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"

	"github.com/gin-gonic/gin"
)

// 发布评课请求数据
type evaluationPublishRequest struct {
	CourseId            string  `json:"course_id" binding:"required"` // FIX 加上 binding
	CourseName          string  `json:"course_name"`
	Rate                uint8   `json:"rate"`
	AttendanceCheckType uint8   `json:"attendance_check_type"`
	ExamCheckType       uint8   `json:"exam_check_type"`
	Content             string  `json:"content"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Tags                []uint8 `json:"tags"`
}

type evaluationPublishResponse struct {
	EvaluationId uint32 `json:"evaluation_id"`
}

type evaluationPublishAPIResponse struct {
	Code    int32
	Message string
	Data    interface{}
}

// Publish ...
// @Summary 发布评课
// @Accept  json
// @Produce  json
// @Param evaluationPublishRequest body comment.evaluationPublishRequest true "评课数据"
// @Success 200 {array} evaluationPublishResponse
// @Router /evaluation [post]
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
		Tags:                model.TagArrayToStr(data.Tags),
		IsAnonymous:         data.IsAnonymous,
		IsValid:             true,
		Time:                strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err := evaluation.New(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// 更新数据库中课程的评分信息
	if err := model.UpdateCourseRateByEvaluation(evaluation.Id, data.Rate); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &evaluationPublishResponse{EvaluationId: evaluation.Id})
}
