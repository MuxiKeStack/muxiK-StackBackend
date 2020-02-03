package grade

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

type GetGradeResponse struct {
	HasLicence bool    `json:"has_licence"`
	TotalScore float32 `json:"total_score"` // 总成绩均分
	UsualScore float32 `json:"usual_score"` // 平时均分
	SampleSize uint32  `json:"sample_size"` // 样本数
	Section1   float32 `json:"section_1"`   // 成绩区间1,85以上所占的比例
	Section2   float32 `json:"section_2"`   // 区间2，70-85所占的比例
	Section3   float32 `json:"section_3"`   // 区间3，85以下所占的比例
}

// @Tags grade
// @Summary 获取成绩
// @Param token header string true "token"
// @Param course_id query string true "课程hash id"
// @Success 200 {object} grade.GetGradeResponse
// @Router /grade/ [get]
func Get(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	// 检查该用户是否有查看成绩的许可
	if ok, err := model.UserHasLicence(userId); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	} else if !ok {
		// 无查看成绩许可，未加入计划
		handler.SendResponse(c, nil, &GetGradeResponse{HasLicence: false})
		return
	}

	courseId, ok := c.GetQuery("course_id")
	if !ok {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "No course_id")
		return
	}

	course, err := model.GetHistoryCourseByHashId(courseId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	var result = GetGradeResponse{
		HasLicence: true,
		TotalScore: course.TotalGrade,
		UsualScore: course.UsualGrade,
		SampleSize: course.GradeSampleSize,
	}

	// 样本数不为0，计算各区间的百分比
	if course.GradeSampleSize != 0 {
		result.Section1 = float32(course.GradeSection1 / course.GradeSampleSize)
		result.Section2 = float32(course.GradeSection2 / course.GradeSampleSize)
		result.Section3 = float32(course.GradeSection3 / course.GradeSampleSize)
	}

	handler.SendResponse(c, nil, result)
}
