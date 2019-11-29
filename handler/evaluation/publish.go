package evaluation

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// 发布评课的请求数据
type evaluationPublishRequest struct {
	CourseId            string  `json:"course_id" binding:"required"`
	CourseName          string  `json:"course_name" binding:"required"`
	Rate                float32 `json:"rate" binding:"-"`
	AttendanceCheckType uint8   `json:"attendance_check_type" binding:"-"` // 经常点名/偶尔点名/签到点名，标识为 1/2/3
	ExamCheckType       uint8   `json:"exam_check_type" binding:"-"`       // 无考核/闭卷考试/开卷考试/论文考核，标识为 1/2/3/4
	Content             string  `json:"content" binding:"-"`
	IsAnonymous         bool    `json:"is_anonymous" binding:"-"`
	Tags                []uint8 `json:"tags" binding:"-"`
}

type evaluationPublishResponse struct {
	EvaluationId uint32 `json:"evaluation_id"`
}

// Publish a new  course evaluation.
// @Summary 发布评课
// @Tags evaluation
// @Param token header string true "token"
// @Param data body evaluation.evaluationPublishRequest true "data"
// @Success 200 {object} evaluation.evaluationPublishResponse
// @Router /evaluation/ [post]
func Publish(c *gin.Context) {
	log.Info("Evaluation Publish function is called.")

	var data evaluationPublishRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
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
		CommentNum:          0,
		Tags:                service.TagArrayToStr(data.Tags),
		IsAnonymous:         data.IsAnonymous,
		IsValid:             true,
		Time:                util.GetCurrentTime(),
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
