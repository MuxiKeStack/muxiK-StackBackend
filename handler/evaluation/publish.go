package evaluation

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/MuxiKeStack/muxiK-StackBackend/util/securityCheck"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Request data of publishing a new evaluation
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
	var data evaluationPublishRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	// Check whether the course exists
	if ok := model.IsCourseExisting(data.CourseId); !ok {
		handler.SendBadRequest(c, errno.ErrCourseExisting, nil, "")
		return
	}

	// Check whether user has evaluated the course
	if ok := model.HasEvaluated(userId, data.CourseId); ok {
		handler.SendBadRequest(c, errno.ErrHasEvaluated, nil, "")
		return
	}

	// Words are limited to 500
	if len(data.Content) > 500 {
		handler.SendBadRequest(c, errno.ErrWordLimitation, nil, "Evaluation's content is limited to 500.")
		return
	}

	// 小程序内容安全检测
	if err := securityCheck.MsgSecCheck(data.Content); err != nil {
		log.Errorf(err, "QQ security check msg(%s) error", data.Content)
		handler.SendBadRequest(c, errno.ErrSecurityCheck, nil, "evaluation content violation")
		return
	}
	log.Info("QQ security check OK")

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

	// Update rate info of the evaluated course
	if err := model.UpdateCourseRateByEvaluation(evaluation.CourseId, data.Rate); err != nil {
		log.Info("UpdateCourseRateByEvaluation function error")
		handler.SendError(c, errno.ErrUpdateCourseInfo, nil, err.Error())
		return
	}

	// Update the tag amount of the course
	if err := service.NewTagsAfterPublishing(&data.Tags, data.CourseId); err != nil {
		log.Info("NewTagsAfterPublishing function error")
		handler.SendError(c, errno.ErrUpdateCourseInfo, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &evaluationPublishResponse{EvaluationId: evaluation.Id})
}
